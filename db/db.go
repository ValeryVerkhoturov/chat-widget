package db

import (
	"context"
	"github.com/ValeryVerkhoturov/chat/config"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	MongoClient *mongo.Client
	Database    *mongo.Database
	UsersColl   *mongo.Collection
)

func InitDB() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	if MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(config.MongoDBUri)); err != nil {
		log.Fatal(err)
	}
	if err = MongoClient.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}
	Database = MongoClient.Database("chat")
	UsersColl = Database.Collection("users")
	return ctx, cancel
}

func ConvertInsertOneResultToId(result *mongo.InsertOneResult) (string, bool) {
	objectId, ok := result.InsertedID.(primitive.ObjectID)
	return objectId.Hex(), ok
}
