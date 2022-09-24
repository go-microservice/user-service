// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/redis"
	"github.com/go-microservice/user-service/internal/cache"
	"github.com/go-microservice/user-service/internal/model"
	"github.com/go-microservice/user-service/internal/repository"
	"github.com/go-microservice/user-service/internal/server"
	"github.com/go-microservice/user-service/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

func InitApp(cfg *app.Config, config *app.ServerConfig) (*app.App, func(), error) {
	db, cleanup, err := model.Init()
	if err != nil {
		return nil, nil, err
	}
	client, cleanup2, err := redis.Init()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	userCache := cache.NewUserCache(client)
	userRepo := repository.NewUser(db, userCache)
	userServiceServer := service.NewUserServiceServer(userRepo)
	grpcServer := server.NewGRPCServer(config, userServiceServer)
	appApp := newApp(cfg, grpcServer)
	return appApp, func() {
		cleanup2()
		cleanup()
	}, nil
}
