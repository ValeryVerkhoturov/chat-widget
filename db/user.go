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
	Source    UserSourceType     `bson:"source"`
	CreatedAt time.Time          `bson:"created_at"`
	IsAgent   bool               `bson:"is_agent"`
}

type UserSourceType string

const (
	EmbeddedChat UserSourceType = "embeddedChat"
	Telegram     UserSourceType = "telegram"
)

func (u *User) FindOneById() (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	err := DB.UsersCollection.FindOne(ctx, bson.D{bson.E{Key: "_id", Value: u.ID}}).Decode(&user)
	return &user, err
}

func (u *User) Find() (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return DB.UsersCollection.Find(ctx, u)
}

func (u *User) InsertOne() (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if u.ID.IsZero() {
		u.ID = primitive.NewObjectID()
	}
	u.CreatedAt = time.Now()
	return DB.UsersCollection.InsertOne(ctx, u)
}
