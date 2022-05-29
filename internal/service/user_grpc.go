package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-microservice/user-service/internal/cache"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-microservice/user-service/internal/tasks"

	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/auth"
	"github.com/go-eagle/eagle/pkg/errcode"
	pb "github.com/go-microservice/user-service/api/user/v1"
	"github.com/go-microservice/user-service/internal/ecode"
	"github.com/go-microservice/user-service/internal/model"
	"github.com/go-microservice/user-service/internal/repository"
	"github.com/go-microservice/user-service/internal/types"
	"github.com/google/wire"
	"github.com/jinzhu/copier"
)

var (
	_ pb.UserServiceServer = (*UserServiceServer)(nil)
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewUserServiceServer)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer

	repo repository.UserRepo
}

func NewUserServiceServer(repo repository.UserRepo) *UserServiceServer {
	return &UserServiceServer{
		repo: repo,
	}
}

func (s *UserServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	var userBase *model.UserModel
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
	_, err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	// send welcome email
	task, err := tasks.NewEmailWelcomeTask(user.Username)
	if err == nil {
		_, _ = tasks.GetClient().Enqueue(task)
	}

	return &pb.RegisterReply{
		Username: req.Username,
	}, nil
}

func newUser(username, email, password string) (*model.UserModel, error) {
	return &model.UserModel{
		Username:  username,
		Email:     email,
		Password:  password,
		Status:    int32(pb.StatusType_NORMAL),
		CreatedAt: time.Now().Unix(),
	}, nil
}

func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	if len(req.Email) == 0 && len(req.Username) == 0 {
		return nil, ecode.ErrInvalidArgument.Status(req).Err()
	}

	// get user base info
	var (
		user *model.UserModel
		err  error
	)
	if req.Email != "" {
		user, err = s.repo.GetUserByEmail(ctx, req.Email)
		if err != nil {
			return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
				"msg": err.Error(),
			})).Status(req).Err()
		}
	}
	if user == nil && len(req.Username) > 0 {
		user, err = s.repo.GetUserByUsername(ctx, req.Username)
		if err != nil {
			return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
				"msg": err.Error(),
			})).Status(req).Err()
		}
	}
	if user != nil && user.ID == 0 {
		return nil, ecode.ErrUserNotFound.Status(req).Err()
	}

	if !auth.ComparePasswords(user.Password, req.Password) {
		return nil, ecode.ErrPasswordIncorrect.Status(req).Err()
	}

	// Sign the json web token.
	payload := map[string]interface{}{"user_id": user.ID, "username": user.Username}
	token, err := app.Sign(ctx, payload, app.Conf.JwtSecret, int64(cache.UserTokenExpireTime))
	if err != nil {
		return nil, ecode.ErrToken.Status(req).Err()
	}

	// record token to redis
	err = cache.NewUserTokenCache().SetUserTokenCache(ctx, user.ID, token, cache.UserTokenExpireTime)
	if err != nil {
		return nil, ecode.ErrToken.Status(req).Err()
	}

	return &pb.LoginReply{
		Id:    user.ID,
		Token: token,
	}, nil
}

func (s *UserServiceServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*emptypb.Empty, error) {
	c := cache.NewUserTokenCache()
	// check token
	token, err := c.GetUserTokenCache(ctx, req.Id)
	if err != nil {
		return nil, ecode.ErrToken.Status(req).Err()
	}
	if token != req.Token {
		return nil, ecode.ErrAccessDenied.Status(req).Err()
	}

	// delete token from cache
	err = c.DelUserTokenCache(ctx, req.GetId())
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	// NOTE: don't set to nil
	return &emptypb.Empty{}, nil
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
	id, err := s.repo.CreateUser(ctx, user)
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
	user := &model.UserModel{
		Username:  req.Username,
		Email:     req.Email,
		UpdatedAt: time.Now().Unix(),
	}
	err := s.repo.UpdateUser(ctx, req.UserId, user)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	return &pb.UpdateUserReply{}, nil
}

func (s *UserServiceServer) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordReply, error) {
	user := &model.UserModel{
		Password:  req.Password,
		UpdatedAt: time.Now().Unix(),
	}
	err := s.repo.UpdateUser(ctx, req.UserId, user)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	return &pb.UpdatePasswordReply{}, nil
}
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	user, err := s.repo.GetUser(ctx, req.Id)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound.WithDetails(errcode.NewDetails(map[string]interface{}{
				"msg": err.Error(),
			})).Status(req).Err()
		}
		return nil, err
	}
	u, err := convertUser(user)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserReply{
		User: u,
	}, nil
}
func (s *UserServiceServer) BatchGetUsers(ctx context.Context, req *pb.BatchGetUsersRequest) (*pb.BatchGetUsersReply, error) {
	// check rpc request if canceled
	if ctx.Err() == context.Canceled {
		return nil, ecode.ErrCanceled.Status(req).Err()
	}

	if len(req.GetIds()) == 0 {
		return nil, errors.New("ids is empty")
	}
	var (
		ids   []int64
		users []*pb.User
	)
	ids = req.GetIds()

	// user base
	userBases, err := s.repo.BatchGetUser(ctx, ids)
	if err != nil {
		return nil, ecode.ErrInternalError.Status(req).Err()
	}
	userMap := make(map[int64]*model.UserModel, 0)
	for _, val := range userBases {
		userMap[val.ID] = val
	}

	// compose data
	for _, id := range ids {
		user, ok := userMap[id]
		if !ok {
			continue
		}
		u, err := convertUser(user)
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

func convertUser(u *model.UserModel) (*pb.User, error) {
	if u == nil {
		return nil, nil
	}
	user := &types.User{
		Id:        u.ID,
		Username:  u.Username,
		Phone:     u.Phone,
		Email:     u.Email,
		LoginAt:   u.LoginAt,
		Status:    u.Status,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Gender:    u.Gender,
		Birthday:  u.Birthday,
		Bio:       u.Bio,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	// copy to pb.user
	pbUser := &pb.User{}
	err := copier.Copy(pbUser, &user)
	if err != nil {
		return nil, err
	}
	return pbUser, nil
}
