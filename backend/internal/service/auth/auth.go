package auth

import (
	"context"
	"crypto/rsa"
	"io"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	media "github.com/kust1q/Zapp/backend/internal/storage/objects"
)

type authCache interface {
	Add(ctx context.Context, dataType, data string) error
	Exists(ctx context.Context, dataType, data string) (bool, error)
}

type authStorage interface {
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
}

type mediaStorage interface {
	Upload(ctx context.Context, file io.Reader, mediaType media.MediaType, filename string) (string, string, error)
	Remove(ctx context.Context, objectPath string) error
}

type tokenStorage interface {
	Store(ctx context.Context, token, userID string, ttl time.Duration) error
}

type AuthServiceConfig struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type authService struct {
	cfg     AuthServiceConfig
	storage authStorage
	cache   authCache
	media   mediaStorage
	tokens  tokenStorage
}

func NewAuthService(cfg AuthServiceConfig, storage authStorage, cache authCache, media mediaStorage, tokens tokenStorage) *authService {
	return &authService{
		cfg:     cfg,
		storage: storage,
		cache:   cache,
		media:   media,
		tokens:  tokens,
	}
}
