package service

import (
	"context"

	pb "github.com/go-microservice/account-service/api/user/v1"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	return &pb.RegisterReply{}, nil
}
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	return &pb.LoginReply{}, nil
}
func (s *UserService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileReply, error) {
	return &pb.UpdateProfileReply{}, nil
}
func (s *UserService) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordReply, error) {
	return &pb.UpdatePasswordReply{}, nil
}
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	return &pb.GetUserReply{}, nil
}
func (s *UserService) BatchGetUsers(ctx context.Context, req *pb.BatchGetUsersRequest) (*pb.BatchGetUsersReply, error) {
	return &pb.BatchGetUsersReply{}, nil
}
