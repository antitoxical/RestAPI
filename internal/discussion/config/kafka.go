package config

import (
	"github.com/IBM/sarama"
	"log"
	"time"
)

const (
	InTopic  = "message-in"
	OutTopic = "message-out"
)

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers []string
	Config  *sarama.Config
}

// NewKafkaConfig creates a new Kafka configuration
func NewKafkaConfig(brokers []string) *KafkaConfig {
	config := sarama.NewConfig()

	// Producer config
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 100 * time.Millisecond
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner

	// Consumer config
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Network config for Docker
	config.Net.MaxOpenRequests = 1
	config.Net.DialTimeout = 10 * time.Second
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second

	// Version config
	config.Version = sarama.V2_8_0_0

	return &KafkaConfig{
		Brokers: brokers,
		Config:  config,
	}
}

// CreateTopics creates Kafka topics if they don't exist
func (k *KafkaConfig) CreateTopics() error {
	admin, err := sarama.NewClusterAdmin(k.Brokers, k.Config)
	if err != nil {
		return err
	}
	defer admin.Close()

	topics := []string{InTopic, OutTopic}
	for _, topic := range topics {
		err := admin.CreateTopic(topic, &sarama.TopicDetail{
			NumPartitions:     3,
			ReplicationFactor: 1, // For local development
		}, false)
		if err != nil && err != sarama.ErrTopicAlreadyExists {
			log.Printf("Error creating topic %s: %v", topic, err)
			return err
		}
	}

	return nil
}
