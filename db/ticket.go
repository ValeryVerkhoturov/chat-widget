package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Ticket struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	Status    string             `bson:"status"`   // open, pending, closed
	Owner     primitive.ObjectID `bson:"owner"`    // Reference to User
	Assignee  primitive.ObjectID `bson:"assignee"` // Reference to User (Agent)
}
