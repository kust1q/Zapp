package auth

import (
	"time"

	"github.com/kust1q/Zapp/backend/internal/config"
)

type authService struct {
	cfg    config.AuthServiceConfig
	db     dataStorage
	tokens tokenStorage
	media  mediaService
	search searchRepository
}

func NewAuthService(cfg config.AuthServiceConfig, db dataStorage, media mediaService, tokens tokenStorage, search searchRepository) *authService {
	return &authService{
		cfg:    cfg,
		db:     db,
		media:  media,
		tokens: tokens,
		search: search,
	}
}

func (s *authService) GetRefreshTTL() time.Duration {
	return s.cfg.RefreshTTL
}
