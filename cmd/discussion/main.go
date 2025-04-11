package main

import (
	"RESTAPI/internal/discussion/api"
	"RESTAPI/internal/discussion/repository"
	"RESTAPI/internal/discussion/service"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize Cassandra connection
	cluster := gocql.NewCluster("localhost")
	cluster.Port = 9042
	cluster.Keyspace = "distcomp"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	defer session.Close()

	// Drop the existing table if it exists
	err = session.Query(`DROP TABLE IF EXISTS tbl_message`).Exec()
	if err != nil {
		log.Printf("Warning: Failed to drop table: %v", err)
	}

	// Create the message table with id as primary key
	err = session.Query(`
		CREATE TABLE IF NOT EXISTS tbl_message (
			id bigint,
			newsid bigint,
			country text,
			content text,
			PRIMARY KEY (id)
		)`).Exec()
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Create index on newsid
	err = session.Query(`
		CREATE INDEX IF NOT EXISTS idx_newsid ON tbl_message (newsid)
	`).Exec()
	if err != nil {
		log.Printf("Warning: Failed to create index: %v", err)
	}

	// Wait for schema changes to propagate
	time.Sleep(2 * time.Second)

	// Initialize components
	messageRepo := repository.NewCassandraMessageRepository(session)
	messageService := service.NewMessageService(messageRepo)
	handler := api.NewHandler(messageService)

	// Set up router
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Start server
	addr := "localhost:24130"
	fmt.Printf("Discussion service starting on %s\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
