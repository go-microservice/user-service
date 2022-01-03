package service

import (
	"context"
	"fmt"

	"account-temp/internal/repository"
)

type IGreeterService interface {
	SayHi(ctx context.Context, name string) (string, error)
}

type greeterService struct {
	repo repository.Repository
}

var _ IGreeterService = (*greeterService)(nil)

func newGreeterService(svc *service) *greeterService {
	return &greeterService{repo: svc.repo}
}

func (s *greeterService) SayHi(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("Hi %s", name), nil
}
