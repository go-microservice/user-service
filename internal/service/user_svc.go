package service

import (
	"context"

	"github.com/go-microservice/account-service/internal/model"

	"github.com/go-microservice/account-service/internal/repository"
)

type IUserService interface {
	Register(ctx context.Context, data *model.UserInfoModel) (string, error)
	Login(ctx context.Context, data *model.UserProfileModel) error
}

type userService struct {
	repo repository.Repository
}

var _ IUserService = (*userService)(nil)

func newUserService(svc *service) *userService {
	return &userService{repo: svc.repo}
}

func (s *userService) Register(ctx context.Context, data *model.UserInfoModel) (string, error) {
	return "", nil
}

func (s *userService) Login(ctx context.Context, data *model.UserProfileModel) error {
	return nil
}
