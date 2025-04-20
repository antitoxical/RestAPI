package config

import (
	"time"
)

// Config holds all configuration for the service
type Config struct {
	Kafka  *KafkaConfig
	DB     *DBConfig
	Server *ServerConfig
}

// DBConfig holds database configuration
type DBConfig struct {
	Hosts       []string
	Keyspace    string
	Consistency string
	Timeout     time.Duration
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port string
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		Kafka: NewKafkaConfig([]string{"localhost:9092"}),
		DB: &DBConfig{
			Hosts:       []string{"localhost"},
			Keyspace:    "distcomp",
			Consistency: "quorum",
			Timeout:     5 * time.Second,
		},
		Server: &ServerConfig{
			Port: ":24130",
		},
	}
}
