package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/auth"
	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/google/wire"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"

	pb "github.com/go-microservice/user-service/api/micro/user/v1"
	"github.com/go-microservice/user-service/internal/ecode"
	"github.com/go-microservice/user-service/internal/model"
	"github.com/go-microservice/user-service/internal/repository"
	"github.com/go-microservice/user-service/internal/types"
)

var (
	_ pb.UserServiceServer = (*UserServiceServer)(nil)
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewUserServiceServer)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer

	repo        repository.UserBaseRepo
	profileRepo repository.UserProfileRepo
}

func NewUserServiceServer(br repository.UserBaseRepo, pr repository.UserProfileRepo) *UserServiceServer {
	return &UserServiceServer{
		repo:        br,
		profileRepo: pr,
	}
}

func (s *UserServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	var userBase *model.UserBaseModel
	// check user is exist
	userBase, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	userBase, err = s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	if userBase != nil && userBase.ID > 0 {
		return nil, ecode.ErrUserIsExist.Status(req).Err()
	}

	// gen a hash password
	pwd, err := auth.HashAndSalt(req.Password)
	if err != nil {
		return nil, errcode.ErrEncrypt
	}

	// create a new user
	user, err := newUser(req.Username, req.Email, pwd)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	_, err = s.repo.CreateUserBase(ctx, user)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	return &pb.RegisterReply{
		Username: req.Username,
	}, nil
}

func newUser(username, email, password string) (*model.UserBaseModel, error) {
	return &model.UserBaseModel{
		Username:  username,
		Email:     email,
		Password:  password,
		Status:    int32(pb.StatusType_NORMAL),
		CreatedAt: time.Now().Unix(),
	}, nil
}

func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	if req.Email == "" && req.Username == "" {
		return nil, ecode.ErrInvalidArgument.Status(req).Err()
	}

	// get user base info
	var (
		userBase *model.UserBaseModel
		err      error
	)
	if req.Email != "" {
		userBase, err = s.repo.GetUserByEmail(ctx, req.Email)
		if err != nil {
			return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
				"msg": err.Error(),
			})).Status(req).Err()
		}
	}
	if userBase == nil && req.Username != "" {
		userBase, err = s.repo.GetUserByUsername(ctx, req.Username)
		if err != nil {
			return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
				"msg": err.Error(),
			})).Status(req).Err()
		}
	}
	if userBase != nil && userBase.ID == 0 {
		return nil, ecode.ErrUserNotFound.Status(req).Err()
	}

	if !auth.ComparePasswords(userBase.Password, req.Password) {
		return nil, ecode.ErrPasswordIncorrect.Status(req).Err()
	}

	// Sign the json web token.
	payload := map[string]interface{}{"user_id": userBase.ID, "username": userBase.Username}
	token, err := app.Sign(ctx, payload, app.Conf.JwtSecret, 86400)
	if err != nil {
		return nil, ecode.ErrToken.Status(req).Err()
	}

	return &pb.LoginReply{
		Token: token,
	}, nil
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	// gen a hash password
	pwd, err := auth.HashAndSalt(req.Password)
	if err != nil {
		return nil, errcode.ErrEncrypt
	}

	// create a new user
	user, err := newUser(req.Username, req.Email, pwd)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	id, err := s.repo.CreateUserBase(ctx, user)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	return &pb.CreateUserReply{
		Id:       id,
		Username: req.Username,
	}, nil
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	user := &model.UserBaseModel{
		Username:  req.Username,
		Email:     req.Email,
		UpdatedAt: time.Now().Unix(),
	}
	err := s.repo.UpdateUserBase(ctx, req.UserId, user)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	return &pb.UpdateUserReply{}, nil
}

func (s *UserServiceServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileReply, error) {
	return &pb.UpdateProfileReply{}, nil
}
func (s *UserServiceServer) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordReply, error) {
	user := &model.UserBaseModel{
		Password:  req.Password,
		UpdatedAt: time.Now().Unix(),
	}
	err := s.repo.UpdateUserBase(ctx, req.UserId, user)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	return &pb.UpdatePasswordReply{}, nil
}
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	userBase, err := s.repo.GetUserBase(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	userProfile, err := s.profileRepo.GetUserProfile(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	u, err := convertUser(userBase, userProfile)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserReply{
		User: u,
	}, nil
}
func (s *UserServiceServer) BatchGetUsers(ctx context.Context, req *pb.BatchGetUsersRequest) (*pb.BatchGetUsersReply, error) {
	if len(req.GetIds()) == 0 {
		return nil, errors.New("ids is empty")
	}

	idsStr := strings.Split(req.Ids, ",")
	var (
		ids   []int64
		users []*pb.User
	)
	for _, v := range idsStr {
		ids = append(ids, cast.ToInt64(v))
	}

	// user base
	userBases, err := s.repo.BatchGetUserBase(ctx, ids)
	if err != nil {
		return nil, ecode.ErrInternalError.Status(req).Err()
	}
	userBaseMap := make(map[int64]*model.UserBaseModel, 0)
	for _, val := range userBases {
		userBaseMap[val.ID] = val
	}
	// user profile
	userProfiles, err := s.profileRepo.BatchGetUserProfile(ctx, ids)
	if err != nil {
		return nil, err
	}
	userProfileMap := make(map[int64]*model.UserProfileModel, 0)
	for _, val := range userProfiles {
		userProfileMap[val.ID] = val
	}

	// compose data
	for _, id := range ids {
		userBase, ok := userBaseMap[id]
		if !ok {
			continue
		}
		userProfile, ok := userProfileMap[id]
		if !ok {
			userProfile = &model.UserProfileModel{}
		}
		u, err := convertUser(userBase, userProfile)
		if err != nil {
			// record log
			continue
		}
		users = append(users, u)
	}

	return &pb.BatchGetUsersReply{
		Users: users,
	}, nil
}

func convertUser(userBase *model.UserBaseModel, userProfile *model.UserProfileModel) (*pb.User, error) {
	user := &types.User{
		Id:        userBase.ID,
		Username:  userBase.Username,
		Phone:     userBase.Phone,
		Email:     userBase.Email,
		LoginAt:   userBase.LoginAt,
		Status:    userBase.Status,
		Nickname:  userProfile.Nickname,
		Avatar:    userProfile.Avatar,
		Gender:    userProfile.Gender,
		Birthday:  userProfile.Birthday,
		Bio:       userProfile.Bio,
		CreatedAt: userBase.CreatedAt,
		UpdatedAt: userBase.UpdatedAt,
	}

	// copy to pb.user
	u := &pb.User{}
	err := copier.Copy(u, &user)
	if err != nil {
		return nil, err
	}
	return u, nil
}
