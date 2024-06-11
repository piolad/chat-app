package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "message-data-centre/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const (
	port                    = ":50051"
	mongoDBConnectionString = "mongodb://adminUser:adminPassword@message-data-centre-db:27017"
	mongoDBName             = "message-db"
)

// server is used to implement service.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements service.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	s.SaveMessage(ctx, in)
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// SaveMessage implements service.GreeterServer
func (s *server) SaveMessage(ctx context.Context, in *pb.HelloRequest) (*pb.HelloRequest, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoDBConnectionString))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(mongoDBName).Collection("messages")

	// Create a BSON document from the input message
	messageDocument := bson.M{
		"message":   in.GetName(), // Accessing the "name" field from HelloRequest
		"timestamp": time.Now(),
		"status":    "unread",
	}

	// Insert the document into the collection
	_, err = collection.InsertOne(context.Background(), messageDocument)
	if err != nil {
		return nil, err
	}

	return in, nil
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
		{Key: "message", Value: 1},
		{Key: "timestamp", Value: -1},
		{Key: "status", Value: 1},
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

func main() {
	if err := ensureCollectionExists_Messages(); err != nil {
		log.Fatalf("failed to ensure collection exists: %v", err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
