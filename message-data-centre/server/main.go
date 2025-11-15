package main

import (
	"context"
	"log"
	"net"

	pb "message-data-centre/proto"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const (
	port                    = ":50051"
	mongoDBConnectionString = "mongodb://adminUser:adminPassword@message-data-centre-db:27017"
	mongoDBName             = "message-db"
)

// server is used to implement service.MessageServiceServer
type server struct {
	pb.UnimplementedMessageServiceServer
	mongoClient *mongo.Client
}

func main() {
	ctx := context.Background()
	log.Printf("Starting message-data-centre...")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDBConnectionString))
	if err != nil {
		log.Fatalf("Failed to connect to mogno: %v", err)
	}
	s := &server{mongoClient: client}

	if err := s.ensureCollectionExists_Messages(ctx); err != nil {
		log.Fatalf("failed to ensure collection exists: %v", err)
	}

	if err := s.ensureCollectionExists_Conversations(ctx); err != nil {
		log.Fatalf("failed to ensure conversation collection exists: %v", err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterMessageServiceServer(grpcServer, s)

	log.Printf("Server listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
