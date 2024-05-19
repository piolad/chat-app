package main

import (
	"context"
	"log"
	"net"

	"github.com/piolad/chat-app/message-data-centre/protos"

	"google.golang.org/grpc"
)

// server is used to implement messagecentre.MessageCentreServer
type server struct{}

// CreateChat implements messagecentre.MessageCentreServer
func (s *server) CreateChat(ctx context.Context, req *protos.ChatInfo) (*protos.ChatCreationResponse, error) {
	log.Printf("Received request to create chat with name: %s, description: %s", req.GetName(), req.GetDescription())

	// Here you can implement the logic to create a chat and generate a chat_id
	chatID := "chat123" // Example chat ID

	// Return the chat ID in the response
	return &protos.ChatCreationResponse{
		ChatId: chatID,
	}, nil
}

func main() {
	log.Println("Starting message data centre")

	// Create a TCP listener
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a gRPC server
	s := grpc.NewServer()

	// Register the service implementation to the gRPC server
	protos.RegisterMessageCentreServer(s, &server{})

	// Start the gRPC server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
