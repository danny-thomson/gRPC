package main

import (
	"context"
	"log"
	"time"

	pb "grpc-example/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = "localhost:50051"
)

type TodoTask struct {
	Name        string
	Description string
	Done        bool
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect : %v", err)
	}

	defer conn.Close()

	client := pb.NewTodoServiceClient(conn)

	defer cancel()

	todos := []TodoTask{
		{Name: "Code review", Description: "Review new feature code", Done: false},
		{Name: "Make YouTube Video", Description: "Start Go for beginners series", Done: false},
		{Name: "Buy groceries", Description: "Buy tomatoes, onions, mangos", Done: false},
	}

	for _, todo := range todos {
		res, err := client.CreateTodo(ctx, &pb.NewTodo{Name: todo.Name, Description: todo.Description, Done: todo.Done})

		if err != nil {
			log.Fatalf("could not create user: %v", err)
		}

		log.Printf(`
           ID : %s
           Name : %s
           Description : %s
           Done : %v,
       `, res.GetId(), res.GetName(), res.GetDescription(), res.GetDone())
	}
}
