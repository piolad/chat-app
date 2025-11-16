package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "message-data-centre/proto"
	"message-data-centre/server/config"
	"message-data-centre/server/grpcserver"
	"message-data-centre/server/service"
	"message-data-centre/server/storage"
)

func main() {
	cfg := config.Load()

	client, err := storage.NewClient(cfg.MongoConnectionString)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(nil)

	if err := storage.EnsureIndexes(client, cfg.MongoDBName); err != nil {
		log.Fatalf("failed to ensure indexes: %v", err)
	}

	msgStore := storage.NewMessageStore(client, cfg.MongoDBName)
	convStore := storage.NewConversationStore(client, cfg.MongoDBName)

	msgService := service.NewMessageService(msgStore, convStore)
	grpcServer := grpcserver.NewServer(msgService)

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMessageServiceServer(s, grpcServer)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())
}
