package repository

import (
	"context"
	"fmt"
	"log"
	"github.com/gocql/gocql"
	"RESTAPI/internal/discussion/model"
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

// Create inserts a new message into Cassandra
func (r *CassandraMessageRepository) Create(ctx context.Context, message *model.Message) error {
	log.Printf("Creating message with ID: %d, NewsID: %d, Content: %s", message.ID, message.NewsID, message.Content)
	
	// If ID is not set, generate a new one
	if message.ID == 0 {
		// Get the maximum ID from the table
		var maxID int64
		err := r.session.Query(`
			SELECT MAX(id) FROM tbl_message`).
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

	// Check if message with this ID already exists
	existing, err := r.FindByID(ctx, message.ID)
	if err != nil {
		return fmt.Errorf("failed to check if message exists: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("message with ID %d already exists", message.ID)
	}

	// Use INSERT IF NOT EXISTS to prevent race conditions
	applied, err := r.session.Query(`
		INSERT INTO tbl_message (id, newsid, country, content)
		VALUES (?, ?, ?, ?)
		IF NOT EXISTS`,
		message.ID, message.NewsID, message.Country, message.Content).
		WithContext(ctx).ScanCAS()
	if err != nil {
		log.Printf("Error creating message: %v", err)
		return fmt.Errorf("failed to create message: %v", err)
	}
	if !applied {
		return fmt.Errorf("message with ID %d already exists", message.ID)
	}

	// Verify the message was created
	created, err := r.FindByID(ctx, message.ID)
	if err != nil {
		return fmt.Errorf("failed to verify message creation: %v", err)
	}
	if created == nil || created.NewsID != message.NewsID {
		log.Printf("Created message verification failed. Expected NewsID: %d, Got: %v", 
			message.NewsID, created)
		return fmt.Errorf("message creation verification failed")
	}

	return nil
}

// FindByID retrieves a message by its ID
func (r *CassandraMessageRepository) FindByID(ctx context.Context, id int64) (*model.Message, error) {
	log.Printf("Finding message by ID: %d", id)
	
	var message model.Message
	err := r.session.Query(`
		SELECT id, newsid, country, content
		FROM tbl_message
		WHERE id = ?`,
		id).WithContext(ctx).Scan(&message.ID, &message.NewsID, &message.Country, &message.Content)

	if err != nil {
		if err == gocql.ErrNotFound {
			log.Printf("Message with ID %d not found", id)
			return nil, nil
		}
		log.Printf("Error finding message by ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to retrieve message: %v", err)
	}

	log.Printf("Found message with ID: %d, NewsID: %d, Content: %s", 
		message.ID, message.NewsID, message.Content)
	return &message, nil
}

// FindByNewsID retrieves all messages for a specific news item
func (r *CassandraMessageRepository) FindByNewsID(ctx context.Context, newsID int64) ([]*model.Message, error) {
	log.Printf("Finding messages by NewsID: %d", newsID)
	var messages []*model.Message
	iter := r.session.Query(`
		SELECT id, newsid, country, content
		FROM tbl_message
		WHERE newsid = ?`,
		newsID).WithContext(ctx).Iter()

	var message model.Message
	for iter.Scan(&message.ID, &message.NewsID, &message.Country, &message.Content) {
		messages = append(messages, &model.Message{
			ID:      message.ID,
			Country: message.Country,
			NewsID:  message.NewsID,
			Content: message.Content,
		})
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
	log.Printf("Updating message with ID: %d, NewsID: %d, Content: %s", 
		message.ID, message.NewsID, message.Content)
	
	// Ensure newsId is set
	if message.NewsID == 0 {
		return fmt.Errorf("newsId is required")
	}

	// Use UPDATE IF EXISTS to handle race conditions
	applied, err := r.session.Query(`
		UPDATE tbl_message
		SET newsid = ?, country = ?, content = ?
		WHERE id = ?
		IF EXISTS`,
		message.NewsID, message.Country, message.Content, message.ID).
		WithContext(ctx).ScanCAS()
	if err != nil {
		log.Printf("Error updating message: %v", err)
		return fmt.Errorf("failed to update message: %v", err)
	}
	if !applied {
		// If message doesn't exist, create it
		log.Printf("Message with ID %d not found, creating new one", message.ID)
		return r.Create(ctx, message)
	}

	// Verify the update
	updated, err := r.FindByID(ctx, message.ID)
	if err != nil {
		return fmt.Errorf("failed to verify message update: %v", err)
	}
	if updated == nil || updated.NewsID != message.NewsID {
		log.Printf("Update verification failed. Expected NewsID: %d, Got: %v", 
			message.NewsID, updated)
		return fmt.Errorf("message update verification failed")
	}

	return nil
}

// Delete removes a message by its ID
func (r *CassandraMessageRepository) Delete(ctx context.Context, id int64) error {
	log.Printf("Deleting message with ID: %d", id)
	err := r.session.Query(`
		DELETE FROM tbl_message
		WHERE id = ?`,
		id).WithContext(ctx).Exec()
	if err != nil {
		log.Printf("Error deleting message: %v", err)
		return fmt.Errorf("failed to delete message: %v", err)
	}
	return nil
}

// FindAll retrieves all messages
func (r *CassandraMessageRepository) FindAll(ctx context.Context) ([]*model.Message, error) {
	log.Printf("Finding all messages")
	var messages []*model.Message
	
	iter := r.session.Query(`
		SELECT id, newsid, country, content
		FROM tbl_message`).WithContext(ctx).Iter()

	var message model.Message
	for iter.Scan(&message.ID, &message.NewsID, &message.Country, &message.Content) {
		// Create a new message for each row to avoid pointer issues
		newMessage := &model.Message{
			ID:      message.ID,
			NewsID:  message.NewsID,
			Country: message.Country,
			Content: message.Content,
		}
		messages = append(messages, newMessage)
		log.Printf("Found message: ID=%d, NewsID=%d, Content=%s", 
			newMessage.ID, newMessage.NewsID, newMessage.Content)
	}

	if err := iter.Close(); err != nil {
		log.Printf("Error closing iterator: %v", err)
		return nil, fmt.Errorf("failed to retrieve messages: %v", err)
	}

	log.Printf("Found %d messages", len(messages))
	return messages, nil
} 