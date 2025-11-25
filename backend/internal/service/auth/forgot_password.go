package auth

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

func (s *authService) ForgotPassword(ctx context.Context, req *entity.ForgotPassword) (*entity.Recovery, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.db.GetUserByEmail(ctx, req.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if err := s.tokens.CloseAllSessions(ctx, strconv.Itoa(user.ID)); err != nil {
		return nil, fmt.Errorf("failed to close sessions: %w", err)
	}

	recoveryToken := s.generateRecoveryToken()
	if err = s.tokens.StoreRecovery(ctx, recoveryToken, strconv.Itoa(user.ID), s.cfg.RecoveryTTL); err != nil {
		return nil, fmt.Errorf("failed to store recovery token: %w", err)
	}

	return &entity.Recovery{
		Recovery: recoveryToken,
	}, nil
}

func (s *authService) generateRecoveryToken() string {
	newUUID := uuid.New()
	token := newUUID.String()
	return token
}
