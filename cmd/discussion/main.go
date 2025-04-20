package main

import (
	"RESTAPI/internal/discussion/api"
	"RESTAPI/internal/discussion/config"
	"RESTAPI/internal/discussion/kafka"
	"RESTAPI/internal/discussion/repository"
	"RESTAPI/internal/discussion/service"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Initialize Cassandra connection
	cluster := gocql.NewCluster(cfg.DB.Hosts...)
	cluster.Port = 9042
	cluster.Keyspace = cfg.DB.Keyspace
	cluster.Consistency = gocql.ParseConsistency(cfg.DB.Consistency)
	cluster.Timeout = cfg.DB.Timeout
	cluster.ConnectTimeout = cfg.DB.Timeout
	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
		NumRetries: 3,
		Min:        100 * time.Millisecond,
		Max:        10 * time.Second,
	}

	// Create session without keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	defer session.Close()

	// Create keyspace if not exists with replication settings
	err = session.Query(`
		CREATE KEYSPACE IF NOT EXISTS distcomp
		WITH replication = {
			'class': 'SimpleStrategy',
			'replication_factor': 1
		}`).Exec()
	if err != nil {
		log.Fatalf("Failed to create keyspace: %v", err)
	}

	// Create the message table
	err = session.Query(`
		CREATE TABLE IF NOT EXISTS distcomp.tbl_message (
			id bigint,
			newsid bigint,
			country text,
			content text,
			state text,
			PRIMARY KEY (id)
		)`).Exec()
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Create index on newsid
	err = session.Query(`
		CREATE INDEX IF NOT EXISTS idx_newsid ON distcomp.tbl_message (newsid)
	`).Exec()
	if err != nil {
		log.Printf("Warning: Failed to create index: %v", err)
	}

	// Initialize Kafka
	if err := cfg.Kafka.CreateTopics(); err != nil {
		log.Printf("Warning: Failed to create topics: %v", err)
	}

	// Create Kafka producer
	producer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Initialize components
	messageRepo := repository.NewCassandraMessageRepository(session)
	messageService := service.NewMessageService(messageRepo)

	// Create Kafka consumer
	consumer, err := kafka.NewConsumer(cfg.Kafka, messageService, producer)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Stop()

	// Start consuming messages
	if err := consumer.Start(); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	// Create API handler
	handler := api.NewHandler(messageService)

	// Set up router
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Start server
	fmt.Printf("Discussion service starting on %s\n", cfg.Server.Port)
	if err := http.ListenAndServe(cfg.Server.Port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
