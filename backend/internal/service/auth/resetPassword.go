package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kust1q/Zapp/backend/internal/dto"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

func (s *authService) ResetPassword(ctx context.Context, userID int, req *dto.ResetPasswordRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	req.OldPassword = strings.TrimSpace(req.OldPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)

	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return ErrInvalidPassword
	}

	if err := s.tokens.CloseAllSessions(ctx, strconv.Itoa(userID)); err != nil {
		return fmt.Errorf("failed to close sessions: %w", err)
	}
	newHashPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hashing failed: %w", err)
	}
	return s.storage.UpdateUserPassword(ctx, userID, string(newHashPassword))
}
