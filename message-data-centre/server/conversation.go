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

// FetchLastXMessages retrieves the last X messages between sender and receiver
// FetchLastXMessages retrieves the last X messages between a sender and receiver based on their conversation_id
func (s *server) FetchLastXMessages(ctx context.Context, in *pb.FetchLastXMessagesRequest) (*pb.FetchLastXMessagesResponse, error) {
	// Log the incoming request
	log.Printf("FetchLastXMessages called with: Sender=%s, Receiver=%s, StartingPoint=%d, Count=%d",
	in.GetSender(), in.GetReceiver(), in.GetStartingPoint(), in.GetCount())
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoDBConnectionString))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	log.Println("Connected to MongoDB");

	conversationCollection := client.Database(mongoDBName).Collection("Conversations")

	// First, fetch the conversation_id for the given sender and receiver
	filter := bson.M{
		"$or": []bson.M{
			{"sender": in.GetSender(), "receiver": in.GetReceiver()},
			{"sender": in.GetReceiver(), "receiver": in.GetSender()},
		},
	}

	log.Println("Filter: ", filter)
	var conversation bson.M
	err = conversationCollection.FindOne(context.Background(), filter).Decode(&conversation)
	if err == mongo.ErrNoDocuments {
		return nil, err // No conversation exists between sender and receiver
	} else if err != nil {
		return nil, err // Other error
	}

	conversationID := conversation["_id"].(primitive.ObjectID).Hex()

	log.Println("Conversation ID: ", conversationID)
	// Now fetch messages with this conversation_id
	messageCollection := client.Database(mongoDBName).Collection("messages")
	messageFilter := bson.M{
		"conversation_id": conversationID,
	}

	opts := options.Find().
		SetSkip(int64(in.GetStartingPoint())).
		SetLimit(int64(in.GetCount())).
		SetSort(bson.D{{Key: "timestamp", Value: -1}}) // Sort by timestamp descending

	cursor, err := messageCollection.Find(context.Background(), messageFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	log.Println("Messages found: ", cursor.RemainingBatchLength())

	var messages []*pb.Message
	for cursor.Next(context.Background()) {
		var messageDoc bson.M
		if err := cursor.Decode(&messageDoc); err != nil {
			return nil, err
		}
		// Map BSON document to pb.Message
		message := &pb.Message{
			Sender:    messageDoc["sender"].(string),
			Receiver:  in.GetReceiver(), // The receiver is known from the request
			Message:   messageDoc["message"].(string),
			Timestamp: messageDoc["timestamp"].(string),
		}
		messages = append(messages, message)
	}

	// Count total messages
	totalMessages, err := messageCollection.CountDocuments(context.Background(), messageFilter)
	if err != nil {
		return nil, err
	}

	log.Println("Total messages: ", totalMessages)

	log.Println("end")

	hasMore := (in.GetStartingPoint() + int32(len(messages))) < int32(totalMessages)

	return &pb.FetchLastXMessagesResponse{
		Messages: messages,
		Count:    int32(len(messages)),
		HasMore:  hasMore,
	}, nil
}



// FetchLastXConversations retrieves the last X conversations where the sender is involved
// FetchLastXConversations retrieves the last X conversations where the user is either a sender or a receiver
func (s *server) FetchLastXConversations(ctx context.Context, in *pb.FetchLastXConversationsRequest) (*pb.FetchLastXConversationsResponse, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoDBConnectionString))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	conversationCollection := client.Database(mongoDBName).Collection("Conversations")

	// Define the filter to find conversations where the user is either the sender or the receiver
	filter := bson.M{
		"$or": []bson.M{
			{"sender": in.GetConversationMember()},
			{"receiver": in.GetConversationMember()},
		},
	}

	// Set options for pagination (start_index and count)
	opts := options.Find().
		SetSkip(int64(in.GetStartIndex())).
		SetLimit(int64(in.GetCount())).
		SetSort(bson.D{{Key: "last_timestamp", Value: -1}}) // Sort by last_timestamp descending

	cursor, err := conversationCollection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var pairs []*pb.SenderReceiverPair
	for cursor.Next(context.Background()) {
		var conversationDoc bson.M
		if err := cursor.Decode(&conversationDoc); err != nil {
			return nil, err
		}

		// Determine if the conversationMember is the sender or receiver
		sender := conversationDoc["sender"].(string)
		receiver := conversationDoc["receiver"].(string)
		conversationMember := in.GetConversationMember()

		var pair *pb.SenderReceiverPair
		if conversationMember == sender {
			// The user is the sender
			pair = &pb.SenderReceiverPair{
				Sender:   sender,  // The user is the sender
				Receiver: receiver, // The receiver is the other person
			}
		} else {
			// The user is the receiver
			pair = &pb.SenderReceiverPair{
				Sender:   sender,  // The sender is the other person
				Receiver: receiver, // The user is the receiver
			}
		}

		pairs = append(pairs, pair)
	}

	// Count total conversations involving the user
	totalConversations, err := conversationCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	hasMore := (in.GetStartIndex() + int32(len(pairs))) < int32(totalConversations)

	return &pb.FetchLastXConversationsResponse{
		Pairs:   pairs,
		Count:   int32(len(pairs)),
		HasMore: hasMore,
	}, nil
}

