package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/sirupsen/logrus"
)

func (s *authService) SignOut(ctx context.Context, req *entity.Refresh) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req.Refresh = strings.TrimSpace(req.Refresh)

	if err := s.tokens.RemoveRefresh(ctx, req.Refresh); err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			logrus.Warn("refresh token already deleted or expired")
		}
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	logrus.Info("refresh token successfully removed")
	return nil
}
