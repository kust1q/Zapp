package auth

import (
	"errors"
	"time"

	"github.com/kust1q/Zapp/backend/internal/config"
)

var (
	ErrUsernameAlreadyUsed = errors.New("username already used")
	ErrEmailAlreadyUsed    = errors.New("email already used")
	ErrInvalidInput        = errors.New("invalid input data")
	ErrInvalidCredentials  = errors.New("invalid credential")
	ErrTokenNotFound       = errors.New("refresh token not found")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUserNotFound        = errors.New("user not found")
)

type authService struct {
	cfg    config.AuthServiceConfig
	db     dataStorage
	tokens tokenStorage
	cache  authCache
	media  mediaService
}

func NewAuthService(cfg config.AuthServiceConfig, db dataStorage, cache authCache, media mediaService, tokens tokenStorage) *authService {
	return &authService{
		cfg:    cfg,
		db:     db,
		cache:  cache,
		media:  media,
		tokens: tokens,
	}
}

func (s *authService) GetRefreshTTL() time.Duration {
	return s.cfg.RefreshTTL
}
