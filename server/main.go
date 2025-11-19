package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/KulbhushanBhalerao/grpc-proto-demo-golang/proto"
	"google.golang.org/grpc"
)

// Server implements the GreetingService
type server struct {
	pb.UnimplementedGreetingServiceServer
}

// SayHello implements the simple RPC method
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received request from: %s", req.GetName())

	// Create response
	response := &pb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s! Welcome to gRPC with Go!", req.GetName()),
		Count:   1,
	}

	return response, nil
}

// SayHelloMultiple implements the server streaming RPC method
func (s *server) SayHelloMultiple(req *pb.HelloRequest, stream pb.GreetingService_SayHelloMultipleServer) error {
	log.Printf("Received streaming request from: %s", req.GetName())

	// Send 5 greetings with a delay
	for i := 1; i <= 5; i++ {
		response := &pb.HelloResponse{
			Message: fmt.Sprintf("Hello #%d, %s! Streaming response %d of 5", i, req.GetName(), i),
			Count:   int32(i),
		}

		if err := stream.Send(response); err != nil {
			return err
		}

		log.Printf("Sent streaming response #%d to %s", i, req.GetName())
		time.Sleep(1 * time.Second) // Simulate some processing time
	}

	return nil
}

func main() {
	// Listen on TCP port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register our service implementation
	pb.RegisterGreetingServiceServer(s, &server{})

	log.Printf("✅ gRPC Server is running on port 50051...")
	log.Printf("Waiting for client connections...")

	// Start serving requests
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
