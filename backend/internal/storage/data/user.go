package data

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/storage/postgres"
)

type userStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *userStorage {
	return &userStorage{
		db: db,
	}
}

func (s *userStorage) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	query := fmt.Sprintf("INSERT INTO %s (username, email, password, bio, gen, avatar_url, created_at, is_superuser) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at", postgres.UserTable)

	var id int
	var createdAt time.Time
	err := s.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.Bio,
		user.Gen,
		user.AvatarURL,
		time.Now(),
		user.IsSuperuser,
	).Scan(&id, &createdAt)

	if err != nil {
		return entity.User{}, err
	}

	return entity.User{
		ID:          id,
		Username:    user.Username,
		Email:       user.Email,
		Password:    user.Password,
		Bio:         user.Bio,
		Gen:         user.Gen,
		AvatarURL:   user.AvatarURL,
		CreatedAt:   createdAt,
		IsSuperuser: false,
	}, nil
}

func (s *userStorage) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE email = $1)", postgres.UserTable)
	var user entity.User
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Bio,
		&user.Gen,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.IsSuperuser,
	)

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (s *userStorage) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE username = $1)", postgres.UserTable)
	var user entity.User
	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Bio,
		&user.Gen,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.IsSuperuser,
	)

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
