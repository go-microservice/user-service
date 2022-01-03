package service

import (
	"github.com/go-microservice/account-service/internal/repository"
)

// Svc global var
var Svc Service

// Service define all service
type Service interface {
	Users() IUserService
}

// service struct
type service struct {
	repo repository.Repository
}

// New init service
func New(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Users() IUserService {
	return newUserService(s)
}
