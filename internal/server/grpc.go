package server

import (
	"time"

	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/transport/grpc"

	v1 "github.com/go-microservice/account-service/api/user/v1"
	"github.com/go-microservice/account-service/internal/service"
)

// NewGRPCServer creates a gRPC server
func NewGRPCServer(cfg *app.ServerConfig) *grpc.Server {

	grpcServer := grpc.NewServer(
		grpc.Network("tcp"),
		grpc.Address(":9090"),
		grpc.Timeout(3*time.Second),
	)

	// register biz service
	v1.RegisterUserServiceServer(grpcServer, service.NewUserServiceServer())

	return grpcServer
}
