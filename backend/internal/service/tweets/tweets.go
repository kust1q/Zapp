package tweets

import (
	"context"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
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
	GetCounts(ctx context.Context, tweetID int) (likes, retweets, replies int, err error)
}

type tweetService struct {
	storage tweetStorage
}

func NewTweetService(storage tweetStorage) *tweetService {
	return &tweetService{
		storage: storage,
	}
}
