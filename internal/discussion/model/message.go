package model

type MessageState string

const (
	StatePending MessageState = "PENDING"
	StateApprove MessageState = "APPROVE"
	StateDecline MessageState = "DECLINE"
)

// Message represents a discussion message in the system
type Message struct {
	ID      int64        `json:"id" cql:"id"`
	Country string       `json:"country" cql:"country"`
	NewsID  int64        `json:"newsId" cql:"newsid"`
	Content string       `json:"content" cql:"content"`
	State   MessageState `json:"state"`
}

// MessageTable represents the Cassandra table name for messages
const MessageTable = "tbl_messages"
