package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type message struct {
	Message        string `bson:"message" json:"message"`
	Timestamp      string `bson:"timestamp" json:"timestamp"`
	ConversationID string `bson:"conversation_id" json:"conversation_id"`
	Sender         string `bson:"sender" json:"sender"`
}

// Saves a message to the database from the conversation service
func (s *server) SaveMessage(ctx context.Context, in *message) error {

	collection := s.mongoClient.Database(mongoDBName).Collection("messages")

	// Create a BSON document from the input message
	messageDocument := bson.M{
		"message":         in.Message,
		"timestamp":       in.Timestamp,
		"status":          "unread",
		"conversation_id": in.ConversationID,
		"sender":          in.Sender,
	}

	// Insert the document into the collection
	_, err := collection.InsertOne(context.Background(), messageDocument)
	if err != nil {
		return err
	}

	return nil
}

func (s *server) ensureCollectionExists_Messages(ctx context.Context) error {
	collection := s.mongoClient.Database(mongoDBName).Collection("messages")

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
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return err
	}

	return nil
}

func (s *server) findOrCreateConversation(ctx context.Context, sender string, receiver string) (string, error) {
	collection := s.mongoClient.Database(mongoDBName).Collection("Conversations")

	// Define the filter to check for existing conversation in both directions
	filter := bson.M{
		"$or": []bson.M{
			{"sender": sender, "receiver": receiver},
			{"sender": receiver, "receiver": sender},
		},
	}

	var result bson.M
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		newConversation := bson.M{
			"sender":         sender,
			"receiver":       receiver,
			"last_timestamp": time.Now(),
			"iv_vector":      "",
		}

		insertResult, err := collection.InsertOne(context.Background(), newConversation)
		if err != nil {
			log.Printf("findOrCreateConversation failed: %v", err)
			return "", err
		}

		// Extract the inserted ID
		conversationID := insertResult.InsertedID.(primitive.ObjectID).Hex()
		log.Println("New conversation created between sender:", sender, "and receiver:", receiver)
		return conversationID, nil
	} else if err != nil {

		log.Printf("findOrCreateConversation failed: %v", err)
		return "", err
	}

	// If a document is found, extract the ID and log that the conversation already exists
	conversationID := result["_id"].(primitive.ObjectID).Hex()
	sender = result["sender"].(string)
	receiver = result["receiver"].(string)
	log.Println("Conversation already exists between sender:", sender, "and receiver:", receiver)
	return conversationID, nil
}

func (s *server) ensureCollectionExists_Conversations(ctx context.Context) error {
	collection := s.mongoClient.Database(mongoDBName).Collection("Conversations")

	//Define the index key
	indexKeys := bson.D{
		{Key: "id", Value: 1},
		{Key: "sender", Value: 1},
		{Key: "receiver", Value: 1},
		{Key: "last_timestamp", Value: 1},
		{Key: "iv_vector", Value: 1},
	}

	// Define the index options
	indexOptions := options.Index().SetUnique(true)

	// Create the index model
	indexModel := mongo.IndexModel{
		Keys:    indexKeys,
		Options: indexOptions,
	}

	// Create the index
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return err
	}

	return nil
}
