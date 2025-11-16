package storage

import (
	"context"
	"log"

	"message-data-centre/server/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageStore struct {
	col *mongo.Collection
}

func NewMessageStore(client *mongo.Client, dbName string) *MessageStore {
	col := client.Database(dbName).Collection("messages")
	log.Printf("NewMessageStore init: db=%s col=%s", col.Database().Name(), col.Name())

	return &MessageStore{
		col: col,
	}
}

func (s *MessageStore) Save(ctx context.Context, m *data.Message) error {
	doc := bson.M{
		"message":         m.Message,
		"timestamp":       m.Timestamp,
		"status":          "unread",
		"conversation_id": m.ConversationID,
		"sender":          m.Sender,
	}
	_, err := s.col.InsertOne(ctx, doc)
	return err
}

func (s *MessageStore) FetchByConversation(
	ctx context.Context, conversationID string, start, count int32,
) ([]*data.Message, int32, error) {
	filter := bson.M{"conversation_id": conversationID}

	opts := options.Find().
		SetSkip(int64(start)).
		SetLimit(int64(count)).
		SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cursor, err := s.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var msgs []*data.Message
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, err
		}
		msgs = append(msgs, &data.Message{
			Message:        doc["message"].(string),
			Timestamp:      doc["timestamp"].(string),
			ConversationID: conversationID,
			Sender:         doc["sender"].(string),
		})
	}

	total, err := s.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return msgs, int32(total), nil
}
