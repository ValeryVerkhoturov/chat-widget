package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id"`
	TicketID  primitive.ObjectID `bson:"ticket_id"` // Reference to Ticket
	SenderID  primitive.ObjectID `bson:"sender_id"` // Reference to User
	Content   string             `bson:"content"`
	CreatedAt time.Time          `bson:"created_at"`
	Type      string             `bson:"type"` // text, attachment
}
