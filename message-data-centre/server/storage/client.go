package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(uri string) (*mongo.Client, error) {
	return mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
}
