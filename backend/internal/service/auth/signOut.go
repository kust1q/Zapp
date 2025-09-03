package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kust1q/Zapp/backend/internal/dto"
	"github.com/sirupsen/logrus"
)

func (s *authService) SignOut(ctx context.Context, token *dto.RefreshRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	token.Refresh = strings.TrimSpace(token.Refresh)

	if err := s.tokens.Remove(ctx, token.Refresh); err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			logrus.Debug("refresh token already deleted or expired")
		}
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	logrus.Info("refresh token successfully removed")
	return nil
}
