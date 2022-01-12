package service

import (
	"context"
	"time"

	"github.com/go-eagle/eagle/pkg/errcode"

	"github.com/go-eagle/eagle/pkg/auth"

	"github.com/go-microservice/account-service/internal/ecode"
	"github.com/go-microservice/account-service/internal/model"
	"github.com/go-microservice/account-service/internal/repository"
)

const (
	UserStatusUnknown = iota
	UserStatusNormal
	UserStatusDelete
	UserStatusBan
)

type IUserService interface {
	Register(ctx context.Context, username, email, password string) error
	Login(ctx context.Context, data *model.UserProfileModel) error
}

type userService struct {
	repo repository.Repository
}

var _ IUserService = (*userService)(nil)

func newUserService(svc *service) *userService {
	return &userService{repo: svc.repo}
}

func (s *userService) Register(ctx context.Context, username, email, password string) error {
	var userInfo *model.UserInfoModel
	// check user is exist
	userInfo, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	userInfo, err = s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	if userInfo != nil && userInfo.ID > 0 {
		return ecode.ErrUserIsExist
	}

	// gen a hash password
	pwd, err := auth.HashAndSalt(password)
	if err != nil {
		return errcode.ErrEncrypt
	}

	// if not exist, register a new user
	data := &model.UserInfoModel{
		Username:  username,
		Email:     email,
		Password:  pwd,
		Status:    UserStatusNormal,
		CreatedAt: time.Now().Unix(),
	}
	_, err = s.repo.CreateUserInfo(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) Login(ctx context.Context, data *model.UserProfileModel) error {
	return nil
}
