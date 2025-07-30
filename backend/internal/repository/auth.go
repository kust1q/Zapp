package repository

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kust1q/Zapp/backend/internal/config"
	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/repository/postgres"
	"github.com/minio/minio-go/v7"
)

type authStorage struct {
	db *sqlx.DB
	mc *minio.Client
}

func NewAuthStorage(db *sqlx.DB, mc *minio.Client) *authStorage {
	return &authStorage{
		db: db,
		mc: mc,
	}
}

func (s *authStorage) UploadAvatar(ctx context.Context, avatarBuffer bytes.Buffer, avatarSaveName string) (string, error) {
	minioConfig := config.Get().Minio
	_, err := s.mc.PutObject(
		ctx,
		minioConfig.BucketName,
		avatarSaveName,
		bytes.NewReader(avatarBuffer.Bytes()),
		int64(avatarBuffer.Len()),
		minio.PutObjectOptions{
			ContentType: "image/png",
		})

	if err != nil {
		return "", err
	}

	var protocol string
	if !minioConfig.MinioUseSSL {
		protocol = "http"
	} else {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s/%s/%s", protocol, minioConfig.MinioEndpoint, minioConfig.BucketName, avatarSaveName), nil
}

func (s *authStorage) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
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
