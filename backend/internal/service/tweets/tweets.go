package tweets

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/dto"
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
	GetRepliesToParentTweet(ctx context.Context, parentTweetID int) ([]entity.Tweet, error)
	GetTweetsAndRetweetsByUsername(ctx context.Context, username string) ([]entity.Tweet, error)
	GetCounts(ctx context.Context, tweetID int) (likes, retweets, replies int, err error)
	GetLikes(ctx context.Context, tweetID int) ([]dto.UserLikeResponse, error)

	GetUserByID(ctx context.Context, userID int) (*entity.User, error)
}

type mediaService interface {
	UploadAndAttachTweetMediaTx(ctx context.Context, tweetID int, file io.Reader, filename string, tx *sql.Tx) (*entity.TweetMedia, error)
	GetMediaByTweetID(ctx context.Context, tweetID int) (*entity.TweetMedia, error)
	GetAvatarByUserID(ctx context.Context, userID int) (*entity.Avatar, error)
}

var (
	ErrTweetNotFound      = errors.New("tweet not found")
	ErrUserNotFound       = errors.New("user not found ")
	ErrUnauthorizedUpdate = errors.New("user is not authorized to update this tweet")
)

type tweetService struct {
	storage tweetStorage
	media   mediaService
}

func NewTweetService(storage tweetStorage, media mediaService) *tweetService {
	return &tweetService{
		storage: storage,
		media:   media,
	}
}

func (s *tweetService) tweetResponseByTweet(ctx context.Context, tweet *entity.Tweet) (*dto.TweetResponse, error) {
	var m dto.TweetMediaResponse
	media, err := s.media.GetMediaByTweetID(ctx, tweet.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			m = dto.TweetMediaResponse{}
		} else {
			return &dto.TweetResponse{}, fmt.Errorf("failed to get tweet by id: %w", err)
		}
	} else {
		m = dto.TweetMediaResponse{
			ID:        media.ID,
			TweetID:   media.TweetID,
			MediaURL:  media.MediaURL,
			MimeType:  media.MimeType,
			SizeBytes: media.SizeBytes,
		}
	}

	author, err := s.storage.GetUserByID(ctx, tweet.UserID)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("failed to get tweet author")
	}

	avatar, err := s.media.GetAvatarByUserID(ctx, tweet.UserID)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("failed to get user avatar")
	}

	return &dto.TweetResponse{
		ID:            tweet.ID,
		Content:       tweet.Content,
		CreatedAt:     tweet.CreatedAt,
		UpdatedAt:     tweet.UpdatedAt,
		ParentTweetID: &tweet.ParentTweetID,
		Media:         m,
		Author: dto.UserResponse{
			ID:       author.ID,
			Username: author.Username,
			Avatar: dto.Avatar{
				MediaURL:  avatar.MediaURL,
				MimeType:  avatar.MimeType,
				SizeBytes: avatar.SizeBytes,
			},
		},
	}, nil
}

func (s *tweetService) tweetResponseWithCountersByTweet(ctx context.Context, tweet *entity.Tweet) (*dto.TweetResponseWithCounters, error) {
	var m dto.TweetMediaResponse
	media, err := s.media.GetMediaByTweetID(ctx, tweet.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			m = dto.TweetMediaResponse{}
		} else {
			return &dto.TweetResponseWithCounters{}, fmt.Errorf("failed to get tweet by id: %w", err)
		}
	} else {
		m = dto.TweetMediaResponse{
			ID:        media.ID,
			TweetID:   media.TweetID,
			MediaURL:  media.MediaURL,
			MimeType:  media.MimeType,
			SizeBytes: media.SizeBytes,
		}
	}

	author, err := s.storage.GetUserByID(ctx, tweet.UserID)
	if err != nil {
		return &dto.TweetResponseWithCounters{}, fmt.Errorf("failed to get tweet author")
	}

	avatar, err := s.media.GetAvatarByUserID(ctx, tweet.UserID)
	if err != nil {
		return &dto.TweetResponseWithCounters{}, fmt.Errorf("failed to get user avatar")
	}

	likes, retweets, replyCount, err := s.storage.GetCounts(ctx, tweet.ID)
	if err != nil {
		return &dto.TweetResponseWithCounters{}, fmt.Errorf("failed to get user counters")
	}

	return &dto.TweetResponseWithCounters{
		TweetResponse: dto.TweetResponse{
			ID:            tweet.ID,
			Content:       tweet.Content,
			CreatedAt:     tweet.CreatedAt,
			UpdatedAt:     tweet.UpdatedAt,
			ParentTweetID: &tweet.ParentTweetID,
			Media:         m,
			Author: dto.UserResponse{
				ID:       author.ID,
				Username: author.Username,
				Avatar: dto.Avatar{
					MediaURL:  avatar.MediaURL,
					MimeType:  avatar.MimeType,
					SizeBytes: avatar.SizeBytes,
				},
			},
		},
		ReplyCount:   replyCount,
		RetweetCount: retweets,
		LikeCount:    likes,
	}, nil
}
