package tweets

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

var (
	ErrTweetNotFound      = errors.New("tweet not found")
	ErrUserNotFound       = errors.New("user not found ")
	ErrUnauthorizedUpdate = errors.New("user is not authorized to update this tweet")
)

type tweetService struct {
	db     tweetStorage
	media  mediaService
	search searchRepository
}

func NewTweetService(db tweetStorage, media mediaService, search searchRepository) *tweetService {
	return &tweetService{
		db:     db,
		media:  media,
		search: search,
	}
}

func (s *tweetService) BuildEntityTweetToResponse(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error) {
	author, err := s.db.GetUserByID(ctx, tweet.Author.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tweet author: %w", err)
	}

	avatarUrl, err := s.media.GetAvatarUrlByUserID(ctx, tweet.Author.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user avatar: %w", err)
	}

	counts, err := s.db.GetCounts(ctx, tweet.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get counters: %w", err)
	}

	return &entity.Tweet{
		ID:            tweet.ID,
		ParentTweetID: tweet.ParentTweetID,
		Content:       tweet.Content,
		CreatedAt:     tweet.CreatedAt,
		UpdatedAt:     tweet.UpdatedAt,
		MediaUrl:      tweet.MediaUrl,
		Author: &entity.SmallUser{
			ID:        author.ID,
			Username:  author.Username,
			AvatarURL: avatarUrl,
		},
		Counters: counts,
	}, nil
}
