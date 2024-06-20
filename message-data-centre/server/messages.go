package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Saves a message to the database from the conversation service
func (s *server) SaveMessage(ctx context.Context, in *message) error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoDBConnectionString))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(mongoDBName).Collection("messages")

	// Create a BSON document from the input message
	messageDocument := bson.M{
		"message":         in.Message,
		"timestamp":       in.Timestamp,
		"status":          "unread",
		"conversation_id": in.ConversationID,
		"sender":          in.Sender,
	}

	// Insert the document into the collection
	_, err = collection.InsertOne(context.Background(), messageDocument)
	if err != nil {
		return err
	}

	return nil
}

func ensureCollectionExists_Messages() error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoDBConnectionString))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(mongoDBName).Collection("messages")

	// Define the index keys
	indexKeys := bson.D{
		{Key: "timestamp", Value: -1},
		{Key: "message", Value: 1},
		{Key: "status", Value: 1},
		{Key: "conversation_id", Value: 1},
		{Key: "sender", Value: 1},
	}

	// Define the index options
	indexOptions := options.Index().SetUnique(true)

	// Create the index model
	indexModel := mongo.IndexModel{
		Keys:    indexKeys,
		Options: indexOptions,
	}

	// Create the index
	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return err
	}

	return nil
}
