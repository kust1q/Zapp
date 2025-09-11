package tweets

import (
	"context"
	"errors"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/dto"
)

type tweetStorage interface {
	CreateTweet(ctx context.Context, tweet *entity.Tweet) (entity.Tweet, error)
	GetTweetById(ctx context.Context, tweetID int) (entity.Tweet, error)
	UpdateTweet(ctx context.Context, tweet *entity.Tweet) (entity.Tweet, error)
	DeleteTweet(ctx context.Context, userID, tweetID int) error
	LikeTweet(ctx context.Context, userID, tweetID int) error
	UnLikeTweet(ctx context.Context, userID, tweetID int) error
	Retweet(ctx context.Context, userID, tweetID int, createdAt time.Time) error
	DeleteRetweet(ctx context.Context, userID, retweetID int) error
	GetRepliesToParentTweet(ctx context.Context, parentTweetID int) ([]entity.Tweet, error)
	GetTweetsByUsername(ctx context.Context, username string) ([]entity.Tweet, error)
	GetCounts(ctx context.Context, tweetID int) (likes, retweets, replies int, err error)
	GetLikes(ctx context.Context, tweetID int) ([]dto.UserLikeResponse, error)
}

var (
	ErrTweetNotFound      = errors.New("tweet not found")
	ErrUserNotFound       = errors.New("user not found ")
	ErrUnauthorizedUpdate = errors.New("user is not authorized to update this tweet")
)

type tweetService struct {
	storage tweetStorage
}

func NewTweetService(storage tweetStorage) *tweetService {
	return &tweetService{
		storage: storage,
	}
}
