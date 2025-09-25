package user

import (
	"context"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/dto"
)

func (s *userService) FollowToUser(ctx context.Context, followerID, followingID int) (*dto.FollowResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if followerID == followingID {
		return &dto.FollowResponse{}, fmt.Errorf("impossible to subscribe to yourself")
	}
	return s.storage.FollowToUser(ctx, followerID, followingID, time.Now())
}

func (s *userService) UnfollowUser(ctx context.Context, followerID, followingID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.storage.UnfollowUser(ctx, followerID, followingID)
}

func (s *userService) GetFollowers(ctx context.Context, username string) ([]dto.SmallUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	followersIDs, err := s.storage.GetFollowersIds(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers ids: %w", err)
	}
	users := make([]dto.SmallUserResponse, 0, len(followersIDs))
	for _, id := range followersIDs {
		user, err := s.storage.GetUserByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get user by id: %w", err)
		}
		avatarURL, err := s.media.GetAvatarUrlByUserID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get avatar")
		}
		users = append(users, dto.SmallUserResponse{
			ID:        user.ID,
			Username:  user.Username,
			AvatarURL: avatarURL,
		})
	}
	return users, nil
}

func (s *userService) GetFollowings(ctx context.Context, username string) ([]dto.SmallUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	followingsIDs, err := s.storage.GetFollowingsIds(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers ids: %w", err)
	}
	users := make([]dto.SmallUserResponse, 0, len(followingsIDs))
	for _, id := range followingsIDs {
		user, err := s.storage.GetUserByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get user by id: %w", err)
		}
		avatarUrl, err := s.media.GetAvatarUrlByUserID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get avatar")
		}
		users = append(users, dto.SmallUserResponse{
			ID:        user.ID,
			Username:  user.Username,
			AvatarURL: avatarUrl,
		})
	}
	return users, nil
}
