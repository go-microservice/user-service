package service

import (
	"context"
	"fmt"

	"github.com/go-microservice/account-service/internal/repository"
)

type IUserService interface {
	SayHi(ctx context.Context, name string) (string, error)
}

type userService struct {
	repo repository.Repository
}

var _ IUserService = (*userService)(nil)

func newUserService(svc *service) *userService {
	return &userService{repo: svc.repo}
}

func (s *userService) SayHi(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("Hi %s", name), nil
}
