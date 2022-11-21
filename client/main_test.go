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

// create mock gRPC server for testing
type mockTodoServiceServer struct {
	pb.UnimplementedTodoServiceServer
}

func (*mockTodoServiceServer) CreateTodo(ctx context.Context, todo *pb.NewTodo) (*pb.Todo, error) {

	return &pb.Todo{
		Name:        todo.GetName(),
		Description: todo.GetDescription(),
		Done:        todo.GetDone(),
	}, nil
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterTodoServiceServer(server, &mockTodoServiceServer{})

	// using goroutine because serve function blocks this function
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestDepositClient_Deposit(t *testing.T) {

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewTodoServiceClient(conn)
	resp, err := client.CreateTodo(ctx, &pb.NewTodo{Name: "Test Name", Description: "Testing the description", Done: true})
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}
	log.Printf("Response: %+v", resp)
}
