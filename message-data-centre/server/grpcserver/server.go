package grpcserver

import (
	"context"
	"fmt"
	"log"

	pb "message-data-centre/proto"
	"message-data-centre/server/data"
	"message-data-centre/server/service"
)

type Server struct {
	pb.UnimplementedMessageServiceServer
	svc *service.MessageService
}

func NewServer(svc *service.MessageService) *Server {
	return &Server{svc: svc}
}

func (s *Server) SendMessage(ctx context.Context, in *pb.Message) (*pb.Response, error) {
	log.Printf("SendMessage request received: %v", in)

	msg := &data.Message{
		Message:   in.GetMessage(),
		Timestamp: in.GetTimestamp(),
		Sender:    in.GetSender(),
	}

	if err := s.svc.SendMessage(ctx, msg, in.GetReceiver()); err != nil {
		return nil, err
	}

	respText := fmt.Sprintf(
		"Message sent by %s to %s with message: %s at time %s",
		in.GetSender(), in.GetReceiver(), in.GetMessage(), in.GetTimestamp(),
	)

	log.Printf("SendMessage response: { Message: %v }", respText)
	return &pb.Response{Message: respText}, nil
}

func (s *Server) FetchLastXConversations(ctx context.Context, in *pb.FetchLastXConversationsRequest) (*pb.FetchLastXConversationsResponse, error) {
	log.Printf("FetchLastXConversations request received: %v", in)

	convos, has_more, err := s.svc.FetchLastConversations(ctx, in.GetConversationMember(), in.GetStartIndex(), in.GetCount())

	if err != nil {
		return nil, err
	}

	log.Printf("FetchLastXConversations got convos: %v", convos)

	pairs := []*pb.SenderReceiverPair{}

	for idx := range convos {
		convo := convos[idx]
		log.Printf("FetchLastXConversations message [%d/%d]: {Sender: %s, Receiver: %s} ", idx+1, len(convos), convo.Sender, convo.Receiver)
		pairs = append(pairs, &pb.SenderReceiverPair{Sender: convo.Sender, Receiver: convo.Receiver})
	}

	log.Printf("FetchLastXConversations response: { Pairs: %v, Count: %d, HasMore: %t }", pairs, int32(len(pairs)), has_more)
	return &pb.FetchLastXConversationsResponse{Pairs: pairs, Count: int32(len(pairs)), HasMore: has_more}, nil
}

func (s *Server) FetchLastXMessages(ctx context.Context, in *pb.FetchLastXMessagesRequest) (*pb.FetchLastXMessagesResponse, error) {
	log.Printf("FetchLastXMessages request received")
	msgs, hasMore, err := s.svc.FetchLastMessages(ctx, in.GetSender(), in.GetReceiver(), in.GetStartingPoint(), in.GetCount())

	log.Printf("FetchLastXMessages fetched [%d] messages ", len(msgs))

	if err != nil {
		return nil, err
	}

	convertedMessages := []*pb.Message{}

	for idx := range msgs {
		msg := msgs[idx]
		convertedMessages = append(convertedMessages, &pb.Message{Sender: msg.Sender, Receiver: in.GetReceiver(), Message: msg.Message, Timestamp: msg.Timestamp})
		log.Printf("FetchLastXMessages message [%d/%d]: {Sender: } ", idx+1, len(msgs))
	}

	log.Printf("FetchLastXMessages response: { Messages: %v, HasMore: %t, Count: %d", convertedMessages, hasMore, int32(len(convertedMessages)))
	return &pb.FetchLastXMessagesResponse{Messages: convertedMessages, HasMore: hasMore, Count: int32(len(convertedMessages))}, nil
}
