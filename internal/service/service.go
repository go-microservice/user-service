package service

import (
	"account-temp/internal/repository"
)

// Svc global var
var Svc Service

// Service define all service
type Service interface {
	Greeter() IGreeterService
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

func (s *service) Greeter() IGreeterService {
	return newGreeterService(s)
}
