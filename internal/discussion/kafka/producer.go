package kafka

import (
	"RESTAPI/internal/discussion/config"
	"RESTAPI/internal/discussion/model"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

// Producer handles message production to Kafka
type Producer struct {
	producer sarama.SyncProducer
}

// NewProducer creates a new Kafka producer
func NewProducer(kafkaConfig *config.KafkaConfig) (*Producer, error) {
	producer, err := sarama.NewSyncProducer(kafkaConfig.Brokers, kafkaConfig.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %v", err)
	}

	return &Producer{producer: producer}, nil
}

// SendMessage sends a message to Kafka
func (p *Producer) SendMessage(topic string, message *model.Message) error {
	value, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", message.NewsID)), // Ensure messages from same news go to same partition
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("Message sent to partition %d at offset %d", partition, offset)
	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.producer.Close()
}
