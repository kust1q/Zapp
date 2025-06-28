package repository

import (
	"github.com/jmoiron/sqlx"
)

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

type Repository struct {
	Auth
	Tweet
	User
	Media
	Search
	Feed
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
