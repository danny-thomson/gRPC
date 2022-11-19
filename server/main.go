package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "grpc-example/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedTodoServiceServer
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
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
	fmt.Println("inside main")
}

func (s *server) CreateTodo(ctx context.Context, todo *pb.NewTodo) (*pb.Todo, error) {
	fmt.Println("Desc", todo.GetDescription())
	fmt.Println("Name", todo.GetName())
	fmt.Println("Done", todo.GetDone())

	//response
	return &pb.Todo{Name: "Hello " + todo.GetName(), Description: "Response " +todo.GetDescription()}, nil
}
