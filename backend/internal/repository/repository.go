package repository

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

func NewRepository() *Repository {
	return &Repository{}
}
