//go:build wireinject
// +build wireinject

package main

import (
	"log"

	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/client/consulclient"
	"github.com/go-eagle/eagle/pkg/client/etcdclient"
	"github.com/go-eagle/eagle/pkg/client/nacosclient"
	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/registry"
	"github.com/go-eagle/eagle/pkg/registry/consul"
	"github.com/go-eagle/eagle/pkg/registry/etcd"
	"github.com/go-eagle/eagle/pkg/registry/nacos"
	"github.com/go-eagle/eagle/pkg/transport/grpc"
	"github.com/go-microservice/user-service/internal/cache"
	"github.com/go-microservice/user-service/internal/repository"
	"github.com/go-microservice/user-service/internal/server"
	"github.com/go-microservice/user-service/internal/service"
	"github.com/google/wire"
)

func InitApp(cfg *eagle.Config, config *eagle.ServerConfig) (*eagle.App, func(), error) {
	wire.Build(server.ProviderSet, service.ProviderSet, repository.ProviderSet, cache.ProviderSet, newApp)
	return &eagle.App{}, nil, nil
}

func newApp(cfg *eagle.Config, gs *grpc.Server) *eagle.App {
	return eagle.New(
		eagle.WithName(cfg.Name),
		eagle.WithVersion(cfg.Version),
		eagle.WithLogger(logger.GetLogger()),
		eagle.WithServer(
			// init HTTP server
			server.NewHTTPServer(&cfg.HTTP),
			// init gRPC server
			gs,
		),
		// eagle.WithRegistry(getConsulRegistry()),
	)
}

// create a etcd register
func getEtcdRegistry() registry.Registry {
	client, err := etcdclient.New()
	if err != nil {
		log.Fatal(err)
	}
	return etcd.New(client.Client)
}

// create a consul register
func getConsulRegistry() registry.Registry {
	client, err := consulclient.New()
	if err != nil {
		panic(err)
	}
	return consul.New(client)
}

// create a nacos register
func getNacosRegistry() registry.Registry {
	client, err := nacosclient.New()
	if err != nil {
		panic(err)
	}
	return nacos.New(client)
}
