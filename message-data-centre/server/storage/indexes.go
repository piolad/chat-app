package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureIndexes(client *mongo.Client, dbName string) error {
	if err := ensureMessagesIndexes(client, dbName); err != nil {
		return err
	}
	if err := ensureConversationsIndexes(client, dbName); err != nil {
		return err
	}
	return nil
}

func ensureMessagesIndexes(client *mongo.Client, dbName string) error {
	collection := client.Database(dbName).Collection("messages")

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "timestamp", Value: -1},
			{Key: "message", Value: 1},
			{Key: "status", Value: 1},
			{Key: "conversation_id", Value: 1},
			{Key: "sender", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

func ensureConversationsIndexes(client *mongo.Client, dbName string) error {
	collection := client.Database(dbName).Collection("Conversations")

	// make this unique only if you want one conv per pair
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "sender", Value: 1},
			{Key: "receiver", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}
