package main

import (
	"context"
	pb "github.com/zhuangqi/biubiu/pkg/server/proto/account"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "localhost:9101"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreateUser(ctx, &pb.CreateUserRequest{Username: "admin", Password: "123456"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("id:%d username: %s", r.GetId(), r.GetUsername())
}
