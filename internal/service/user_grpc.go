package service

import (
	"context"
	"time"

	"github.com/jinzhu/copier"

	"github.com/go-microservice/user-service/internal/types"

	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/auth"

	"github.com/go-microservice/user-service/internal/model"
	"github.com/go-microservice/user-service/internal/repository"

	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/go-microservice/user-service/internal/ecode"

	pb "github.com/go-microservice/user-service/api/user/v1"
)

var (
	_ pb.UserServiceServer = (*UserServiceServer)(nil)
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer

	repo repository.Repository
}

func NewUserServiceServer() *UserServiceServer {
	return &UserServiceServer{
		repo: repository.New(model.GetDB()),
	}
}

func (s *UserServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	var userInfo *model.UserInfoModel
	// check user is exist
	userInfo, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	userInfo, err = s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	if userInfo != nil && userInfo.ID > 0 {
		return nil, ecode.ErrUserIsExist.Status(req).Err()
	}

	// gen a hash password
	pwd, err := auth.HashAndSalt(req.Password)
	if err != nil {
		return nil, errcode.ErrEncrypt
	}

	// if not exist, register a new user
	data := &model.UserInfoModel{
		Username:  req.Username,
		Email:     req.Email,
		Password:  pwd,
		Status:    int32(pb.StatusType_NORMAL),
		CreatedAt: time.Now().Unix(),
	}
	_, err = s.repo.CreateUserInfo(ctx, data)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	return &pb.RegisterReply{
		Username: req.Username,
	}, nil
}
func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	// get user base info
	userInfo, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	userInfo, err = s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	if userInfo != nil && userInfo.ID > 0 {
		return nil, ecode.ErrUserIsExist.Status(req).Err()
	}

	// gen a hash password
	pwd, err := auth.HashAndSalt(req.Password)
	if err != nil {
		return nil, ecode.ErrEncrypt.Status(req).Err()
	}
	if pwd != userInfo.Password {
		return nil, ecode.ErrPasswordIncorrect.Status(req).Err()
	}

	// Sign the json web token.
	payload := map[string]interface{}{"user_id": userInfo.ID, "username": userInfo.Username}
	token, err := app.Sign(ctx, payload, app.Conf.JwtSecret, 86400)
	if err != nil {
		return nil, ecode.ErrToken.Status(req).Err()
	}

	return &pb.LoginReply{
		Token: token,
	}, nil
}
func (s *UserServiceServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileReply, error) {
	return &pb.UpdateProfileReply{}, nil
}
func (s *UserServiceServer) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordReply, error) {
	return &pb.UpdatePasswordReply{}, nil
}
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	userInfo, err := s.repo.GetUserInfo(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	userProfile, err := s.repo.GetUserProfile(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	user := &types.User{
		ID:        userInfo.ID,
		Username:  userInfo.Username,
		Phone:     userInfo.Phone,
		Email:     userInfo.Email,
		LoginAt:   userInfo.LoginAt,
		Status:    userInfo.Status,
		Nickname:  userProfile.Nickname,
		Avatar:    userProfile.Avatar,
		Gender:    userProfile.Gender,
		Birthday:  userProfile.Birthday,
		Bio:       userProfile.Bio,
		CreatedAt: userInfo.CreatedAt,
		UpdatedAt: userInfo.UpdatedAt,
	}

	// copy to pb.user
	u := pb.User{}
	err = copier.Copy(&u, &user)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserReply{
		User: &u,
	}, nil
}
func (s *UserServiceServer) BatchGetUsers(ctx context.Context, req *pb.BatchGetUsersRequest) (*pb.BatchGetUsersReply, error) {
	return &pb.BatchGetUsersReply{}, nil
}
