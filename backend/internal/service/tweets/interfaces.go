package tweets

import (
	"context"
	"database/sql"
	"io"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

type tweetStorage interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreateTweetTx(ctx context.Context, tx *sql.Tx, tweet *entity.Tweet) (*entity.Tweet, error)
	CreateTweet(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error)
	GetTweetById(ctx context.Context, tweetID int) (*entity.Tweet, error)
	UpdateTweet(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error)
	DeleteTweet(ctx context.Context, userID, tweetID int) error
	LikeTweet(ctx context.Context, userID, tweetID int) error
	UnLikeTweet(ctx context.Context, userID, tweetID int) error
	Retweet(ctx context.Context, userID, tweetID int, createdAt time.Time) error
	DeleteRetweet(ctx context.Context, userID, retweetID int) error
	GetRepliesToTweet(ctx context.Context, parentTweetID int) ([]entity.Tweet, error)
	GetTweetsAndRetweetsByUsername(ctx context.Context, username string) ([]entity.Tweet, error)
	GetCounts(ctx context.Context, tweetID int) (*entity.Counters, error)
	GetLikes(ctx context.Context, tweetID int) ([]entity.SmallUser, error)

	GetUserByID(ctx context.Context, userID int) (*entity.User, error)
}

type mediaService interface {
	UploadAndAttachTweetMediaTx(ctx context.Context, tweetID int, file io.Reader, filename string, tx *sql.Tx) (string, error)
	GetMediaUrlByTweetID(ctx context.Context, tweetID int) (string, error)
	GetAvatarUrlByUserID(ctx context.Context, userID int) (string, error)
	DeleteTweetMedia(ctx context.Context, tweetID, userID int) error
	GetPresignedURL(ctx context.Context, path string) (string, error)
}

type searchRepository interface {
	IndexTweet(ctx context.Context, tweet *entity.Tweet) error
	DeleteTweet(ctx context.Context, tweetID int) error
}
