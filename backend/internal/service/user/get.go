package user

import (
	"context"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/dto"
)

func (s *userService) GetUserByID(ctx context.Context, userID int) (*dto.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return &dto.UserResponse{}, fmt.Errorf("failed to get user by id: %w", err)
	}

	avatarUrl, err := s.media.GetAvatarUrlByUserID(ctx, user.ID)
	if err != nil {
		return &dto.UserResponse{}, fmt.Errorf("failed to get avatar by user id: %w", err)
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Bio:       user.Bio,
		Gen:       user.Gen,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		AvatarURL: avatarUrl,
	}, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := s.storage.GetUserByUsername(ctx, username)
	if err != nil {
		return &dto.UserResponse{}, fmt.Errorf("failed to get user by id: %w", err)
	}

	avatarURL, err := s.media.GetAvatarUrlByUserID(ctx, user.ID)
	if err != nil {
		return &dto.UserResponse{}, fmt.Errorf("failed to get avatar by user id: %w", err)
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Bio:       user.Bio,
		Gen:       user.Gen,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		AvatarURL: avatarURL,
	}, nil
}
