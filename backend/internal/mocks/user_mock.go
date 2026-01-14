package mocks

import (
	"context"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

type MockUserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	UpdateBio(ctx context.Context, userID int, bio string) error
	FollowUser(ctx context.Context, followerID, followingID int) error
	UnfollowUser(ctx context.Context, followerID, followingID int) error
	GetFollowers(ctx context.Context, username string, limit, offset int) ([]entity.SmallUser, error)
	GetFollowings(ctx context.Context, username string, limit, offset int) ([]entity.SmallUser, error)
	DeleteUser(ctx context.Context, userID int) error
	GetUserProfile(ctx context.Context, username string, limit, offset int) (*entity.UserProfile, error)
}

type MockTweetService interface {
	GetTweetsAndRetweetsByUsername(ctx context.Context, username string, limit, offset int) ([]entity.Tweet, error)
}
