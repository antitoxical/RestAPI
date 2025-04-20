package service

import (
	"RESTAPI/internal/discussion/model"
	"RESTAPI/internal/discussion/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// MessageService handles business logic for messages
type MessageService struct {
	repo   repository.MessageRepository
	client *http.Client
}

// NewMessageService creates a new MessageService
func NewMessageService(repo repository.MessageRepository) *MessageService {
	return &MessageService{
		repo: repo,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *MessageService) checkNewsExists(newsId int64) error {
	resp, err := s.client.Get(fmt.Sprintf("http://localhost:24110/api/v1.0/news/%d", newsId))
	if err != nil {
		return fmt.Errorf("failed to check news existence: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("news with ID %d does not exist", newsId)
	}

	return nil
}

func (s *MessageService) getMessageFromMainService(id int64) (*model.Message, error) {
	resp, err := s.client.Get(fmt.Sprintf("http://localhost:24110/api/v1.0/messages/%d", id))
	if err != nil {
		return nil, fmt.Errorf("failed to get message from main service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("message with ID %d not found in main service", id)
	}

	var message model.Message
	if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
		return nil, fmt.Errorf("failed to decode message: %v", err)
	}

	return &message, nil
}

// CreateMessage creates a new message
func (s *MessageService) CreateMessage(ctx context.Context, message *model.Message) error {
	log.Printf("Creating message with NewsID: %d", message.NewsID)

	// Validate newsId
	if message.NewsID == 0 {
		return fmt.Errorf("newsId is required")
	}

	if err := s.checkNewsExists(message.NewsID); err != nil {
		return err
	}

	// Create the message
	err := s.repo.Create(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to create message: %v", err)
	}

	// Log successful creation
	log.Printf("Successfully created message with ID: %d, NewsID: %d",
		message.ID, message.NewsID)

	return nil
}

// GetMessage retrieves a message by ID
func (s *MessageService) GetMessage(ctx context.Context, id int64) (*model.Message, error) {
	log.Printf("Getting message with ID: %d", id)

	// First try to get from Cassandra
	message, err := s.repo.FindByID(ctx, id)
	if err == nil {
		return message, nil
	}

	// If not found in Cassandra, try to get from main service
	message, err = s.getMessageFromMainService(id)
	if err != nil {
		return nil, err
	}

	// Save to Cassandra for future use
	if err := s.repo.Create(ctx, message); err != nil {
		log.Printf("Warning: Failed to save message to Cassandra: %v", err)
	}

	return message, nil
}

// GetMessagesByNewsID retrieves all messages for a news item
func (s *MessageService) GetMessagesByNewsID(ctx context.Context, newsID int64) ([]*model.Message, error) {
	log.Printf("Getting messages for NewsID: %d", newsID)

	// First try to get from Cassandra
	messages, err := s.repo.FindByNewsID(ctx, newsID)
	if err == nil && len(messages) > 0 {
		return messages, nil
	}

	// If not found in Cassandra, try to get from main service
	resp, err := s.client.Get(fmt.Sprintf("http://localhost:24110/api/v1.0/messages?newsId=%d", newsID))
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from main service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get messages: status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, fmt.Errorf("failed to decode messages: %v", err)
	}

	// Save to Cassandra for future use
	for _, message := range messages {
		if err := s.repo.Create(ctx, message); err != nil {
			log.Printf("Warning: Failed to save message to Cassandra: %v", err)
		}
	}

	return messages, nil
}

// UpdateMessage updates an existing message
func (s *MessageService) UpdateMessage(ctx context.Context, message *model.Message) error {
	log.Printf("Updating message with ID: %d, NewsID: %d", message.ID, message.NewsID)

	// Validate newsId
	if message.NewsID == 0 {
		return fmt.Errorf("newsId is required")
	}

	if err := s.checkNewsExists(message.NewsID); err != nil {
		return err
	}

	// Update the message
	err := s.repo.Update(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to update message: %v", err)
	}

	// Log successful update
	log.Printf("Successfully updated message with ID: %d, NewsID: %d",
		message.ID, message.NewsID)

	return nil
}

// DeleteMessage deletes a message by ID
func (s *MessageService) DeleteMessage(ctx context.Context, id int64) error {
	log.Printf("Deleting message with ID: %d", id)

	return s.repo.Delete(ctx, id)
}

// GetAllMessages retrieves all messages
func (s *MessageService) GetAllMessages(ctx context.Context) ([]*model.Message, error) {
	log.Printf("Getting all messages")

	// First try to get from Cassandra
	messages, err := s.repo.FindAll(ctx)
	if err == nil && len(messages) > 0 {
		return messages, nil
	}

	// If not found in Cassandra, try to get from main service
	resp, err := s.client.Get("http://localhost:24110/api/v1.0/messages")
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from main service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get messages: status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, fmt.Errorf("failed to decode messages: %v", err)
	}

	// Save to Cassandra for future use
	for _, message := range messages {
		if err := s.repo.Create(ctx, message); err != nil {
			log.Printf("Warning: Failed to save message to Cassandra: %v", err)
		}
	}

	return messages, nil
}
