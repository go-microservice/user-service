package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/go-microservice/user-service/api/user/v1"
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

	userClient := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	userReq := &pb.GetUserRequest{
		Id: 1,
	}
	reply, err := userClient.GetUser(ctx, userReq)
	if err != nil {
		log.Fatalf("[rpc] get user err: %v", err)
	}
	fmt.Printf("UserService  GetUser: %+v", reply)
}
