package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Login     string             `bson:"login,omitempty"`
	Email     string             `bson:"email,omitempty"`
	Source    string             `bson:"source"` // embeddedChat, telegram
	CreatedAt time.Time          `bson:"created_at"`
	IsAgent   bool               `bson:"is_agent"`
}

func (u *User) InsertOne() (*mongo.InsertOneResult, error) {
	return UsersColl.InsertOne(context.TODO(), bson.D{
		{"name", u.Login},
		{"email", u.Email},
		{"source", u.Source},
		{"created_at", u.CreatedAt},
		{"is_agent", u.IsAgent},
	})
}
