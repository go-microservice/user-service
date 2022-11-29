package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/go-microservice/user-service/internal/model"

	"github.com/go-microservice/user-service/internal/mocks"

	"github.com/golang/mock/gomock"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/go-microservice/user-service/api/user/v1"
)

const (
	addr    = ""
	bufSize = 1024 * 1024
)

var (
	lis *bufconn.Listener
)

// init create a gRPC server
func initGRPCServer(t *testing.T) {
	lis = bufconn.Listen(bufSize)
	srv := grpc.NewServer()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepo(ctrl)

	pb.RegisterUserServiceServer(srv, &UserServiceServer{
		repo: mockUserRepo,
	})

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("srv.Serve, err: %v", err)
		}
	}()
}

func dialer() func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}

func TestUserServiceServer_GetUser(t *testing.T) {
	testCases := []struct {
		name       string
		id         int64
		res        *pb.GetUserReply
		buildStubs func(mock *mocks.MockUserRepo)
		errCode    codes.Code
		errMsg     string
	}{
		{
			name: "OK",
			id:   1,
			res:  &pb.GetUserReply{User: &pb.User{Id: 1}},
			buildStubs: func(mock *mocks.MockUserRepo) {
				mock.EXPECT().GetUser(gomock.Any(), int64(1)).
					Return(&model.UserModel{ID: 1}, nil).Times(1)
			},
			errCode: codes.OK,
			errMsg:  "",
		},
		{
			name: "NotFound",
			id:   10,
			res:  &pb.GetUserReply{User: &pb.User{Id: 0}},
			buildStubs: func(mock *mocks.MockUserRepo) {
				mock.EXPECT().GetUser(gomock.Any(), int64(10)).
					Return(&model.UserModel{}, nil).Times(1)
			},
			errCode: codes.NotFound,
			errMsg:  "not found",
		},
		{
			name: "InternalError",
			id:   2,
			res:  &pb.GetUserReply{User: &pb.User{Id: 2}},
			buildStubs: func(mock *mocks.MockUserRepo) {
				mock.EXPECT().GetUser(gomock.Any(), int64(2)).
					Return(&model.UserModel{}, errors.New("internal error")).Times(1)
			},
			errCode: codes.Code(10000),
			errMsg:  "Internal error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lis := bufconn.Listen(bufSize)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUserRepo(ctrl)
			tc.buildStubs(mockUserRepo)

			srv := grpc.NewServer()
			pb.RegisterUserServiceServer(srv, &UserServiceServer{
				repo: mockUserRepo,
			})

			go func() {
				if err := srv.Serve(lis); err != nil {
					log.Fatalf("srv.Serve, err: %v", err)
				}
			}()

			dialer := func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}

			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, addr, grpc.WithContextDialer(dialer), grpc.WithInsecure())
			if err != nil {
				log.Fatalf("grpc.DialContext, err: %v", err)
			}
			client := pb.NewUserServiceClient(conn)

			req := &pb.GetUserRequest{Id: tc.id}
			resp, err := client.GetUser(ctx, req)
			if err != nil {
				fmt.Println("~~~~~~~~~~~~~~~~", err)
				if er, ok := status.FromError(err); ok {
					if er.Code() != tc.errCode {
						t.Errorf("error code, expected: %d, received: %d", tc.errCode, er.Code())
					}
					if er.Message() != tc.errMsg {
						t.Errorf("error message, expected: %s, received: %s", tc.errMsg, er.Message())
					}
				}
			}
			if resp != nil {
				if resp.GetUser().Id != tc.res.GetUser().Id {
					t.Errorf("response, expected: %v, received: %v", tc.res.GetUser().Id, resp.GetUser().Id)
				}
			}
		})
	}
}
