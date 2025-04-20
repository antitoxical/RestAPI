package repository

import (
	"RESTAPI/internal/discussion/model"
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"log"
)

// MessageRepository defines the interface for message storage operations
type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) error
	FindAll(ctx context.Context) ([]*model.Message, error)
	FindByID(ctx context.Context, id int64) (*model.Message, error)
	FindByNewsID(ctx context.Context, newsID int64) ([]*model.Message, error)
	Update(ctx context.Context, message *model.Message) error
	Delete(ctx context.Context, id int64) error
}

// CassandraMessageRepository implements MessageRepository using Cassandra
type CassandraMessageRepository struct {
	session *gocql.Session
}

// NewCassandraMessageRepository creates a new CassandraMessageRepository
func NewCassandraMessageRepository(session *gocql.Session) *CassandraMessageRepository {
	return &CassandraMessageRepository{session: session}
}

// validateState checks if the given state is valid
func validateState(state string) error {
	switch state {
	case string(model.StatePending), string(model.StateApprove), string(model.StateDecline):
		return nil
	default:
		return fmt.Errorf("invalid state: %s", state)
	}
}

// Create inserts a new message into Cassandra
func (r *CassandraMessageRepository) Create(ctx context.Context, message *model.Message) error {
	log.Printf("Creating message with ID: %d, NewsID: %d, Content: %s, State: %s",
		message.ID, message.NewsID, message.Content, message.State)

	// If ID is not set, generate a new one
	if message.ID == 0 {
		// Get the maximum ID from the table
		var maxID int64
		err := r.session.Query(`
			SELECT MAX(id) FROM tbl_message
			USING CONSISTENCY QUORUM`).
			WithContext(ctx).Scan(&maxID)
		if err != nil && err != gocql.ErrNotFound {
			log.Printf("Error getting max ID: %v", err)
			return fmt.Errorf("failed to get max ID: %v", err)
		}
		message.ID = maxID + 1
	}

	// Ensure newsId is set
	if message.NewsID == 0 {
		return fmt.Errorf("newsId is required")
	}

	// Set and validate state
	if message.State == "" {
		message.State = model.StatePending
	}
	if err := validateState(string(message.State)); err != nil {
		return err
	}

	// Use INSERT IF NOT EXISTS to prevent race conditions
	applied, err := r.session.Query(`
		INSERT INTO tbl_message (id, newsid, country, content, state)
		VALUES (?, ?, ?, ?, ?)
		IF NOT EXISTS
		USING CONSISTENCY QUORUM`,
		message.ID, message.NewsID, message.Country, message.Content, message.State).
		WithContext(ctx).ScanCAS()
	if err != nil {
		log.Printf("Error creating message: %v", err)
		return fmt.Errorf("failed to create message: %v", err)
	}
	if !applied {
		log.Printf("Message with ID %d already exists", message.ID)
		return fmt.Errorf("message with ID %d already exists", message.ID)
	}

	// Verify the message was created
	created, err := r.FindByID(ctx, message.ID)
	if err != nil {
		log.Printf("Error verifying message creation: %v", err)
		return fmt.Errorf("failed to verify message creation: %v", err)
	}
	if created == nil {
		log.Printf("Message with ID %d was not created", message.ID)
		return fmt.Errorf("message with ID %d was not created", message.ID)
	}

	// Log successful creation
	log.Printf("Successfully created message with ID: %d, NewsID: %d, Content: %s, State: %s",
		message.ID, message.NewsID, message.Content, message.State)

	return nil
}

// FindByID retrieves a message by its ID
func (r *CassandraMessageRepository) FindByID(ctx context.Context, id int64) (*model.Message, error) {
	log.Printf("Finding message by ID: %d", id)

	var message model.Message
	err := r.session.Query(`
		SELECT id, newsid, country, content, state
		FROM tbl_message
		WHERE id = ?
		USING CONSISTENCY QUORUM
	`, id).Scan(&message.ID, &message.NewsID, &message.Country, &message.Content, &message.State)

	if err != nil {
		if err == gocql.ErrNotFound {
			log.Printf("Message with ID %d not found", id)
			return nil, fmt.Errorf("message with ID %d not found", id)
		}
		log.Printf("Error finding message by ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to retrieve message: %v", err)
	}

	log.Printf("Found message with ID: %d, NewsID: %d, Content: %s, State: %s",
		message.ID, message.NewsID, message.Content, message.State)
	return &message, nil
}

// FindByNewsID retrieves all messages for a specific news item
func (r *CassandraMessageRepository) FindByNewsID(ctx context.Context, newsID int64) ([]*model.Message, error) {
	log.Printf("Finding messages by NewsID: %d", newsID)

	iter := r.session.Query(`
		SELECT id, newsid, country, content, state
		FROM tbl_message
		WHERE newsid = ?
		USING CONSISTENCY QUORUM
	`, newsID).Iter()

	var messages []*model.Message
	var message model.Message

	for iter.Scan(&message.ID, &message.NewsID, &message.Country, &message.Content, &message.State) {
		// Create a new message for each iteration to avoid pointer issues
		msg := model.Message{
			ID:      message.ID,
			NewsID:  message.NewsID,
			Country: message.Country,
			Content: message.Content,
			State:   message.State,
		}
		messages = append(messages, &msg)
	}

	if err := iter.Close(); err != nil {
		log.Printf("Error closing iterator: %v", err)
		return nil, fmt.Errorf("failed to retrieve messages: %v", err)
	}

	if len(messages) == 0 {
		log.Printf("No messages found for NewsID: %d", newsID)
		return []*model.Message{}, nil
	}

	log.Printf("Found %d messages for NewsID: %d", len(messages), newsID)
	return messages, nil
}

// Update modifies an existing message
func (r *CassandraMessageRepository) Update(ctx context.Context, message *model.Message) error {
	log.Printf("Updating message with ID: %d, NewsID: %d, Content: %s, State: %s",
		message.ID, message.NewsID, message.Content, message.State)

	// Ensure newsId is set
	if message.NewsID == 0 {
		return fmt.Errorf("newsId is required")
	}

	// Validate state
	if err := validateState(string(message.State)); err != nil {
		return err
	}

	// Use UPDATE IF EXISTS to handle race conditions
	applied, err := r.session.Query(`
		UPDATE tbl_message
		SET newsid = ?, country = ?, content = ?, state = ?
		WHERE id = ?
		IF EXISTS
		USING CONSISTENCY QUORUM`,
		message.NewsID, message.Country, message.Content, message.State, message.ID).
		WithContext(ctx).ScanCAS()
	if err != nil {
		log.Printf("Error updating message: %v", err)
		return fmt.Errorf("failed to update message: %v", err)
	}
	if !applied {
		log.Printf("Message with ID %d not found", message.ID)
		return fmt.Errorf("message with ID %d not found", message.ID)
	}

	// Verify the message was updated
	updated, err := r.FindByID(ctx, message.ID)
	if err != nil {
		log.Printf("Error verifying message update: %v", err)
		return fmt.Errorf("failed to verify message update: %v", err)
	}
	if updated == nil {
		log.Printf("Message with ID %d was not updated", message.ID)
		return fmt.Errorf("message with ID %d was not updated", message.ID)
	}

	// Log successful update
	log.Printf("Successfully updated message with ID: %d, NewsID: %d, Content: %s, State: %s",
		message.ID, message.NewsID, message.Content, message.State)

	return nil
}

// Delete removes a message by its ID
func (r *CassandraMessageRepository) Delete(ctx context.Context, id int64) error {
	log.Printf("Deleting message with ID: %d", id)

	// Check if message exists
	existing, err := r.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check existing message: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("message with ID %d not found", id)
	}

	err = r.session.Query(`
		DELETE FROM tbl_message
		WHERE id = ?
		USING CONSISTENCY QUORUM`,
		id).WithContext(ctx).Exec()
	if err != nil {
		log.Printf("Error deleting message: %v", err)
		return fmt.Errorf("failed to delete message: %v", err)
	}

	// Verify the message was deleted
	deleted, err := r.FindByID(ctx, id)
	if err == nil && deleted != nil {
		log.Printf("Message with ID %d was not deleted", id)
		return fmt.Errorf("message with ID %d was not deleted", id)
	}

	// Log successful deletion
	log.Printf("Successfully deleted message with ID: %d", id)

	return nil
}

// FindAll retrieves all messages
func (r *CassandraMessageRepository) FindAll(ctx context.Context) ([]*model.Message, error) {
	log.Printf("Finding all messages")

	iter := r.session.Query(`
		SELECT id, newsid, country, content, state
		FROM tbl_message
		USING CONSISTENCY QUORUM
	`).Iter()

	var messages []*model.Message
	var message model.Message

	for iter.Scan(&message.ID, &message.NewsID, &message.Country, &message.Content, &message.State) {
		// Create a new message for each iteration to avoid pointer issues
		msg := model.Message{
			ID:      message.ID,
			NewsID:  message.NewsID,
			Country: message.Country,
			Content: message.Content,
			State:   message.State,
		}
		messages = append(messages, &msg)
	}

	if err := iter.Close(); err != nil {
		log.Printf("Error closing iterator: %v", err)
		return nil, fmt.Errorf("failed to retrieve messages: %v", err)
	}

	log.Printf("Found %d messages", len(messages))
	return messages, nil
}
