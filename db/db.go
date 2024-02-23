package db

import (
	"context"
	"errors"
	"github.com/ValeryVerkhoturov/chat/config"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	MongoClient *mongo.Client
	DB          *AppDatabase
)

type AppDatabase struct {
	mongo.Database
	UsersCollection    *mongo.Collection
	TicketsCollection  *mongo.Collection
	MessagesCollection *mongo.Collection
}

func newAppDatabase(db *mongo.Database) *AppDatabase {
	return &AppDatabase{
		Database:           *db,
		UsersCollection:    db.Collection("users"),
		TicketsCollection:  db.Collection("tickets"),
		MessagesCollection: db.Collection("messages"),
	}
}

func CreateDBConnection() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	if MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(config.MongoDBUri)); err != nil {
		log.Fatal(err)
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pingCancel()
	if err = MongoClient.Ping(pingCtx, nil); err != nil {
		log.Fatal(err)
	}
	
	DB = newAppDatabase(MongoClient.Database("chat"))
	DB.UsersCollection = DB.Collection("users")
	DB.TicketsCollection = DB.Collection("tickets")
	DB.MessagesCollection = DB.Collection("messages")

	return ctx, cancel
}

func ConvertInsertOneResultToId(result *mongo.InsertOneResult) (string, bool) {
	objectId, ok := result.InsertedID.(primitive.ObjectID)
	return objectId.Hex(), ok
}

func ConvertValueToObjectId(id interface{}) (primitive.ObjectID, error) {
	switch value := id.(type) {
	case primitive.ObjectID:
		return value, nil
	case string:
		return primitive.ObjectIDFromHex(value)
	default:
		return primitive.NilObjectID, errors.New("object id is not a string")
	}
}
