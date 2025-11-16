package storage

import (
	"context"
	"log"
	"time"

	"message-data-centre/server/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConversationStore struct {
	col *mongo.Collection
}

func NewConversationStore(client *mongo.Client, dbName string) *ConversationStore {
	col := client.Database(dbName).Collection("Conversations")
	log.Printf("NewConversationStore init: db=%s col=%s", col.Database().Name(), col.Name())

	return &ConversationStore{
		col: col,
	}
}

func (s *ConversationStore) GetOrCreate(
	ctx context.Context, sender, receiver string,
) (string, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"sender": sender, "receiver": receiver},
			{"sender": receiver, "receiver": sender},
		},
	}

	var result bson.M
	err := s.col.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		newConversation := bson.M{
			"sender":         sender,
			"receiver":       receiver,
			"last_timestamp": time.Now(),
			"iv_vector":      "",
		}
		insertResult, err := s.col.InsertOne(ctx, newConversation)
		if err != nil {
			return "", err
		}
		return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
	} else if err != nil {
		return "", err
	}

	return result["_id"].(primitive.ObjectID).Hex(), nil
}

func (s *ConversationStore) FetchByMember(
	ctx context.Context, member string, start, count int32,
) ([]*data.Conversation, int32, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"sender": member},
			{"receiver": member},
		},
	}

	opts := options.Find().
		SetSkip(int64(start)).
		SetLimit(int64(count)).
		SetSort(bson.D{{Key: "last_timestamp", Value: -1}})

	log.Printf("FetchByMember using filter: %v", filter)

	log.Printf("FetchByMember using opts skip: %d, %d", *opts.Skip, *opts.Limit)

	cursor, err := s.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var convs []*data.Conversation
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, err
		}
		id := doc["_id"].(primitive.ObjectID).Hex()

		convs = append(convs, &data.Conversation{
			ID:       id,
			Sender:   doc["sender"].(string),
			Receiver: doc["receiver"].(string),
		})
	}

	log.Printf("FetchByMember got convs: %v", convs)

	total, err := s.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return convs, int32(total), nil
}
