package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kust1q/Zapp/backend/internal/dto"
	"github.com/sirupsen/logrus"
)

var (
	ErrTokenNotFound       = errors.New("refresh token not found")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUserNotFound        = errors.New("user not found")
)

func (s *authService) Refresh(ctx context.Context, token *dto.RefreshRequest) (dto.SignInResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	token.Refresh = strings.TrimSpace(token.Refresh)

	userID, err := s.tokens.GetUserIdByRefreshToken(ctx, token.Refresh)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			return dto.SignInResponse{}, ErrInvalidRefreshToken
		}
		return dto.SignInResponse{}, fmt.Errorf("failed to get refresh token: %w", err)
	}

	if err := s.tokens.Remove(ctx, token.Refresh); err != nil {
		logrus.Warnf("failed to delete refresh token: %v", err)
	}

	userIntID, err := strconv.Atoi(userID)
	if err != nil {
		return dto.SignInResponse{}, fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.storage.GetUserByID(ctx, userIntID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.SignInResponse{}, ErrUserNotFound
		}
		return dto.SignInResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	role := "user"
	if user.IsSuperuser {
		role = "admin"
	}
	accessToken, err := s.generateAccessToken(user.ID, user.Email, role)
	if err != nil {
		return dto.SignInResponse{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(ctx, strconv.Itoa(user.ID))
	if err != nil {
		return dto.SignInResponse{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return dto.SignInResponse{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}
