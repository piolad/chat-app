package main

import (
	"context"
	"log"
	pb "message-data-centre/proto"
)

const (
	mongoConversationCollectionName = "Conversations"
	mongoMessageCollectionName      = "Messages"
)

func (s *server) SendMessage(ctx context.Context, in *pb.Message) (*pb.Response, error) {
	log.Printf("Recived message: %v", in.GetMessage())
	return &pb.Response{Message: "Message recived" + in.GetMessage()}, nil
}
