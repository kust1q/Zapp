package user

import (
	"context"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

type dataStorage interface {
	//user
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	UpdateUserBio(ctx context.Context, userID int, bio string) error
	DeleteUser(ctx context.Context, userID int) error
	FollowToUser(ctx context.Context, followerID, followingID int, createdAt time.Time) (*entity.Follow, error)
	UnfollowUser(ctx context.Context, followerID, followingID int) error
	GetFollowersIds(ctx context.Context, username string) ([]int, error)
	GetFollowingsIds(ctx context.Context, username string) ([]int, error)
	//tweets
	GetTweetsAndRetweetsByUsername(ctx context.Context, username string) ([]entity.Tweet, error)
	GetCounts(ctx context.Context, tweetID int) (*entity.Counters, error)
}

type mediaService interface {
	GetAvatarUrlByUserID(ctx context.Context, userID int) (string, error)
	DeleteAvatar(ctx context.Context, userID int) error

	GetMediaUrlByTweetID(ctx context.Context, tweetID int) (string, error)
}

type searchRepository interface {
	DeleteUser(ctx context.Context, userID int) error
	DeleteTweetsByUserID(ctx context.Context, userID int) error
}
