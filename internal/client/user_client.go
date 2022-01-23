package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/go-microservice/user-service/api/micro/user/v1"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	ctx := context.Background()

	userClient := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// get user
	userReq := &pb.GetUserRequest{
		Id: 1,
	}
	reply, err := userClient.GetUser(ctx, userReq)
	if err != nil {
		log.Fatalf("[rpc] get user err: %v", err)
	}
	fmt.Printf("UserService  GetUser: %+v\n", reply)

	// register
	registerReq := &pb.RegisterRequest{
		Username: "test05",
		Email:    "test05@go-eagle.org",
		Password: "123456",
	}
	regReply, err := userClient.Register(ctx, registerReq)
	if err != nil {
		log.Fatalf("[rpc] register err: %v\n", err)
	}
	fmt.Printf("UserService  register resp: %+v\n", regReply)

	// login
	loginReq := &pb.LoginRequest{
		Username: "",
		Email:    "test05@go-eagle.org",
		Password: "123456",
	}
	loginReply, err := userClient.Login(ctx, loginReq)
	if err != nil {
		log.Fatalf("[rpc] login err: %v\n", err)
	}
	fmt.Printf("UserService login resp: %+v\n", loginReply)
}
