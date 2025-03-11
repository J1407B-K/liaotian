package model

import "time"

type Message struct {
	ConversationID string    `bson:"conversation_id"`
	SenderID       string    `bson:"sender_id"`
	ReceiverID     string    `bson:"receiver_id"`
	Content        string    `bson:"content"`
	Timestamp      time.Time `bson:"timestamp"`
	Status         string    `bson:"status"`
}
