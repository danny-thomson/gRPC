package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc-example/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = "localhost:50051"
)

type server struct {
	pb.UnimplementedTodoServiceServer
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterTodoServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) CreateTodo(ctx context.Context, todo *pb.NewTodo) (*pb.Todo, error) {
	fmt.Printf("Name: %s, Description: %s, Done: %t\n", todo.GetName(), todo.GetDescription(), todo.GetDone())

	//response
	return &pb.Todo{
		Name:        todo.GetName(),
		Description: todo.GetDescription(),
		Done:        todo.GetDone(),
	}, nil
}
