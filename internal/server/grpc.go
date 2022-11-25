package server

import (
	"github.com/google/wire"

	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/transport/grpc"

	v1 "github.com/go-microservice/user-service/api/user/v1"
	"github.com/go-microservice/user-service/internal/service"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer)

// NewGRPCServer creates a gRPC server
func NewGRPCServer(cfg *app.ServerConfig, svc *service.UserServiceServer) *grpc.Server {

	grpcServer := grpc.NewServer(
		grpc.Network("tcp"),
		grpc.Address(cfg.Addr),
		grpc.Timeout(cfg.WriteTimeout),
		grpc.EnableLog(),
		grpc.EnableTracing(),
	)

	// register biz service
	v1.RegisterUserServiceServer(grpcServer, svc)

	return grpcServer
}
