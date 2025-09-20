package tweets

import (
	"context"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/dto"
)

func (s *tweetService) CreateTweet(ctx context.Context, userID int, tweet *dto.CreateTweetRequest) (*dto.TweetResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	domainTweet := entity.Tweet{
		UserID:    userID,
		Content:   tweet.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdTweet, err := s.storage.CreateTweet(ctx, &domainTweet)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("tweet creation failed: %w", err)
	}

	author, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("failed to get tweet author")
	}

	avatar, err := s.media.GetAvatarByUserID(ctx, userID)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("failed to get user avatar")
	}

	return &dto.TweetResponse{
		ID:        createdTweet.ID,
		Content:   createdTweet.Content,
		CreatedAt: createdTweet.CreatedAt,
		UpdatedAt: createdTweet.UpdatedAt,
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

func (s *tweetService) CreateTweetWithMedia(ctx context.Context, userID int, tweet *dto.CreateTweetRequest, file *dto.FileData) (*dto.TweetResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	domainTweet := entity.Tweet{
		UserID:    userID,
		Content:   tweet.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdTweet, err := s.storage.CreateTweetTx(ctx, tx, &domainTweet)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("user creation failed: %w", err)
	}

	media, err := s.media.UploadAndAttachTweetMediaTx(ctx, createdTweet.ID, file.File, file.Header.Filename, tx)

	if err := tx.Commit(); err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("commit transaction failed: %w", err)
	}

	author, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("failed to get tweet author")
	}

	avatar, err := s.media.GetAvatarByUserID(ctx, userID)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("failed to get user avatar")
	}

	return &dto.TweetResponse{
		ID:        createdTweet.ID,
		Content:   createdTweet.Content,
		CreatedAt: createdTweet.CreatedAt,
		UpdatedAt: createdTweet.UpdatedAt,
		Media: dto.TweetMediaResponse{
			ID:        media.ID,
			TweetID:   media.TweetID,
			MediaURL:  media.MediaURL,
			MimeType:  media.MimeType,
			SizeBytes: media.SizeBytes,
		},
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
