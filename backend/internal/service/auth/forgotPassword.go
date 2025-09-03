package auth

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kust1q/Zapp/backend/internal/dto"
	"golang.org/x/crypto/bcrypt"
)

func (s *authService) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.SecretAnswer = strings.TrimSpace(req.SecretAnswer)
	req.NewPassword = strings.TrimSpace(req.NewPassword)

	user, err := s.storage.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	securityData, err := s.storage.GetSecurityDataByUserID(ctx, user.ID)
	if err := bcrypt.CompareHashAndPassword([]byte(securityData.Answer), []byte(req.SecretAnswer)); err != nil {
		return fmt.Errorf("invalid request")
	}
	if err := s.tokens.CloseAllSessions(ctx, strconv.Itoa(user.ID)); err != nil {
		return fmt.Errorf("failed to close sessions: %w", err)
	}
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.storage.UpdateUserPassword(ctx, user.ID, string(newPasswordHash))
}
