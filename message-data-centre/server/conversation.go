package main

import (
	"context"
	"log"
	pb "message-data-centre/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *server) SendMessage(ctx context.Context, in *pb.Message) (*pb.Response, error) {
	log.Printf("Recived message: %v", in.GetMessage())
	id, err := s.ensureConversationExists(in.GetSender(), in.GetReceiver())

	if err != nil {
		return nil, err
	}

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

// FetchLastXMessages retrieves the last X messages between sender and receiver
// FetchLastXMessages retrieves the last X messages between a sender and receiver based on their conversation_id
func (s *server) FetchLastXMessages(ctx context.Context, in *pb.FetchLastXMessagesRequest) (*pb.FetchLastXMessagesResponse, error) {
	// Log the incoming request
	log.Printf("FetchLastXMessages called with: Sender=%s, Receiver=%s, StartingPoint=%d, Count=%d",
		in.GetSender(), in.GetReceiver(), in.GetStartingPoint(), in.GetCount())

	conversationCollection := s.mongoClient.Database(mongoDBName).Collection("Conversations")

	// First, fetch the conversation_id for the given sender and receiver
	filter := bson.M{
		"$or": []bson.M{
			{"sender": in.GetSender(), "receiver": in.GetReceiver()},
			{"sender": in.GetReceiver(), "receiver": in.GetSender()},
		},
	}

	log.Println("Filter: ", filter)
	var conversation bson.M
	err := conversationCollection.FindOne(context.Background(), filter).Decode(&conversation)
	if err == mongo.ErrNoDocuments {
		return nil, err // No conversation exists between sender and receiver
	} else if err != nil {
		return nil, err // Other error
	}

	conversationID := conversation["_id"].(primitive.ObjectID).Hex()

	log.Println("Conversation ID: ", conversationID)
	// Now fetch messages with this conversation_id
	messageCollection := s.mongoClient.Database(mongoDBName).Collection("messages")
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
	log.Printf("Data received: ConversationMember=%s, StartIndex=%d, Count=%d", in.GetConversationMember(), in.GetStartIndex(), in.GetCount())

	conversationCollection := s.mongoClient.Database(mongoDBName).Collection("Conversations")

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
				Sender:   sender,   // The user is the sender
				Receiver: receiver, // The receiver is the other person
			}
		} else {
			// The user is the receiver
			pair = &pb.SenderReceiverPair{
				Sender:   sender,   // The sender is the other person
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
