package main

import (
	"context"
	"log"
	pb "message-data-centre/proto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// const (
// 	mongoConversationCollectionName = "Conversations"
// 	mongoMessageCollectionName      = "Messages"
// )

func (s *server) SendMessage(ctx context.Context, in *pb.Message) (*pb.Response, error) {
	log.Printf("Recived message: %v", in.GetMessage())
	id, err := s.ensureConversationExists(in.GetSender(), in.GetReceiver())
	if err != nil {
		return nil, err
	}
	//update timestamp last

	// Create a new message instance with the updated content
	newMessage := &message{
		Message:        in.GetMessage(),
		Timestamp:      in.GetTimestamp(),
		ConversationID: id,
		Sender:         in.GetSender(),
	}

	s.SaveMessage(ctx, newMessage)

	return &pb.Response{Message: "Message send by " + in.GetSender() + " to " + in.GetReceiver() + " with message : " + in.GetMessage() + " at time " + in.GetTimestamp()}, nil
}

func (s *server) ensureConversationExists(sender string, receiver string) (string, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoDBConnectionString))
	if err != nil {
		log.Fatal(err) // error during connection to
		return "", err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(mongoDBName).Collection("Conversations")

	// Define the filter to check for existing conversation in both directions
	filter := bson.M{
		"$or": []bson.M{
			{"sender": sender, "receiver": receiver},
			{"sender": receiver, "receiver": sender},
		},
	}

	var result bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		newConversation := bson.M{
			"sender":         sender,
			"receiver":       receiver,
			"last_timestamp": time.Now(),
			"iv_vector":      "",
		}

		insertResult, err := collection.InsertOne(context.Background(), newConversation)
		if err != nil {
			log.Fatal(err)
			return "", err
		}

		// Extract the inserted ID
		conversationID := insertResult.InsertedID.(primitive.ObjectID).Hex()
		log.Println("New conversation created between sender:", sender, "and receiver:", receiver)
		return conversationID, nil
	} else if err != nil {
		// Handle other potential errors
		log.Fatal(err)
		return "", err
	}

	// If a document is found, extract the ID and log that the conversation already exists
	conversationID := result["_id"].(primitive.ObjectID).Hex()
	sender = result["sender"].(string)
	receiver = result["receiver"].(string)
	log.Println("Conversation already exists between sender:", sender, "and receiver:", receiver)
	return conversationID, nil
}

func ensureCollectionExists_Conversations() error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoDBConnectionString))
	if err != nil {
		log.Fatal(err) //error during connection to database
		return err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(mongoDBName).Collection("Conversations")

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
	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return err
	}

	return nil
}
