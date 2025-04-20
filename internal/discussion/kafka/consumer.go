package kafka

import (
	"RESTAPI/internal/discussion/config"
	"RESTAPI/internal/discussion/model"
	"RESTAPI/internal/discussion/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"strings"
	"sync"
)

// Consumer handles message consumption from Kafka
type Consumer struct {
	consumer       sarama.ConsumerGroup
	messageService *service.MessageService
	producer       *Producer
	stopCh         chan struct{}
	stopOnce       sync.Once
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(kafkaConfig *config.KafkaConfig, messageService *service.MessageService, producer *Producer) (*Consumer, error) {
	group, err := sarama.NewConsumerGroup(kafkaConfig.Brokers, "discussion-group", kafkaConfig.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %v", err)
	}

	return &Consumer{
		consumer:       group,
		messageService: messageService,
		producer:       producer,
		stopCh:         make(chan struct{}),
	}, nil
}

// Start starts consuming messages
func (c *Consumer) Start() error {
	topics := []string{config.InTopic}
	ctx := context.Background()

	go func() {
		for {
			select {
			case <-c.stopCh:
				return
			default:
				if err := c.consumer.Consume(ctx, topics, c); err != nil {
					log.Printf("Error from consumer: %v", err)
				}
			}
		}
	}()

	return nil
}

// Stop stops consuming messages
func (c *Consumer) Stop() {
	c.stopOnce.Do(func() {
		close(c.stopCh)
		if err := c.consumer.Close(); err != nil {
			log.Printf("Error closing consumer: %v", err)
		}
	})
}

// Setup is run at the beginning of a new session
func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session
func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// moderateMessage checks if message contains stop words
func (c *Consumer) moderateMessage(message *model.Message) model.MessageState {
	stopWords := []string{"spam", "abuse", "hate", "violence"}
	content := strings.ToLower(message.Content)

	for _, word := range stopWords {
		if strings.Contains(content, word) {
			return model.StateDecline
		}
	}

	return model.StateApprove
}

// ConsumeClaim processes messages from a partition
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			var msg model.Message
			if err := json.Unmarshal(message.Value, &msg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			// Moderate message
			msg.State = c.moderateMessage(&msg)

			// Save message to database
			if err := c.messageService.CreateMessage(context.Background(), &msg); err != nil {
				log.Printf("Error saving message: %v", err)
				continue
			}

			// Send response to OutTopic
			if err := c.producer.SendMessage(config.OutTopic, &msg); err != nil {
				log.Printf("Error sending response: %v", err)
				continue
			}

			session.MarkMessage(message, "")

		case <-c.stopCh:
			return nil
		}
	}
}
