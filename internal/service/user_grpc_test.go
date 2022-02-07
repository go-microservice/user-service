package service

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/go-microservice/user-service/api/micro/user/v1"
)

const (
	addr    = "localhost:9090"
	bufSize = 1024 * 1024
)

var listener *bufconn.Listener

func initGRPCServerHTTP2() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen, err: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &UserServiceServer{})
	reflection.Register(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to server, err: %v", err)
		}
	}()
}

func initGRPCServerBuffConn() {
	listener = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &UserServiceServer{})
	reflection.Register(s)

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to server, err: %v", err)
		}
	}()
}

func TestUserServiceServer_GetUser(t *testing.T) {
	//initGRPCServerHTTP2()
	initGRPCServerBuffConn()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect, err: %v", err)
	}
	defer conn.Close()

	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetUser(ctx, &pb.GetUserRequest{Id: 1})
	if err != nil {
		log.Fatalf("Could not get user, err: %v", err)
	}
	log.Printf("resp %v", r.User)
}
