package main

import (
	"log"
	"net"

	pb "message-data-centre/proto"

	"google.golang.org/grpc"
)

const (
	port                    = ":50051"
	mongoDBConnectionString = "mongodb://adminUser:adminPassword@message-data-centre-db:27017"
	mongoDBName             = "message-db"
)

// server is used to implement service.GreeterServer.
type server struct {
	pb.UnimplementedMessageServiceServer
}

type message struct {
	Message        string `bson:"message" json:"message"`
	Timestamp      string `bson:"timestamp" json:"timestamp"`
	ConversationID string `bson:"conversation_id" json:"conversation_id"`
	Sender         string `bson:"sender" json:"sender"`
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
	pb.RegisterMessageServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
