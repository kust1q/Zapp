package service

import "github.com/kust1q/Zapp/backend/internal/repository"

type Auth interface {
}

type Tweet interface {
}

type User interface {
}

type Media interface {
}

type Search interface {
}

type Feed interface {
}

type Service struct {
	Auth
	Tweet
	User
	Media
	Search
	Feed
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
