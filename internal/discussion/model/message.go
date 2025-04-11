package model

// Message represents a discussion message in the system
type Message struct {
	ID      int64  `json:"id" cql:"id"`
	Country string `json:"country" cql:"country"`
	NewsID  int64  `json:"newsId" cql:"newsid"`
	Content string `json:"content" cql:"content"`
}

// MessageTable represents the Cassandra table name for messages
const MessageTable = "tbl_messages"
