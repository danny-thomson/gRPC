package main

import (
	"context"
	"log"
	"net"
	"testing"

	pb "grpc-example/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestCreateTodo(t *testing.T) {
	bufSize := 1024 * 1024
	lis := bufconn.Listen(bufSize)

	s := grpc.NewServer()

	pb.RegisterTodoServiceServer(s, &server{})

	// using goroutine because serve function blocks this function
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer conn.Close()

	client := pb.NewTodoServiceClient(conn)
	newTodo := &pb.NewTodo{Name: "Test Name", Description: "Testing the description", Done: true}
	resp, err := client.CreateTodo(ctx, newTodo)
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}

	log.Printf("Response: %+v", resp)

	if newTodo.Name != resp.Name || newTodo.Description != resp.Description || newTodo.Done != resp.Done {
		t.Error("Value is not equal")
	}
}
