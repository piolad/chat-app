package main

import (
	"context"
	"log"
	"net"

	pb "message-data-centre/proto"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement service.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements service.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
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
