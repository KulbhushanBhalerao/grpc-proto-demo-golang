package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/KulbhushanBhalerao/grpc-proto-demo-golang/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := pb.NewGreetingServiceClient(conn)

	log.Println("🚀 gRPC Client started...")
	log.Println("=" + string(make([]byte, 50)) + "=")

	// Example 1: Simple unary RPC call
	fmt.Println("\n📞 Making simple SayHello call...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Alice"})
	if err != nil {
		log.Fatalf("Error calling SayHello: %v", err)
	}

	fmt.Printf("✅ Response: %s\n", response.GetMessage())
	fmt.Printf("   Count: %d\n", response.GetCount())

	// Example 2: Server streaming RPC call
	fmt.Println("\n📡 Making streaming SayHelloMultiple call...")
	stream, err := client.SayHelloMultiple(context.Background(), &pb.HelloRequest{Name: "Bob"})
	if err != nil {
		log.Fatalf("Error calling SayHelloMultiple: %v", err)
	}

	// Receive streaming responses
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// Stream has ended
			fmt.Println("\n✅ Streaming complete!")
			break
		}
		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}

		fmt.Printf("📨 Received: %s (Count: %d)\n", response.GetMessage(), response.GetCount())
	}

	fmt.Println("\n" + string(make([]byte, 50)))
	log.Println("✅ Client finished successfully!")
}
