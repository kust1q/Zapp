package tweets

import (
	"context"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/dto"
)

func (s *tweetService) ReplyToTweet(ctx context.Context, userID, tweetID int, tweet *dto.CreateTweetRequest) (*dto.TweetResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	domainTweet := entity.Tweet{
		UserID:        userID,
		ParentTweetID: tweetID,
		Content:       tweet.Content,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	createdTweet, err := s.storage.CreateTweet(ctx, &domainTweet)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("tweet creation failed: %w", err)
	}

	author, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("failed to get tweet author")
	}

	avatarURL, err := s.media.GetAvatarUrlByUserID(ctx, userID)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("failed to get user avatar")
	}

	return &dto.TweetResponse{
		ID:            createdTweet.ID,
		Content:       createdTweet.Content,
		CreatedAt:     createdTweet.CreatedAt,
		UpdatedAt:     createdTweet.UpdatedAt,
		ParentTweetID: &createdTweet.ParentTweetID,
		Author: dto.SmallUserResponse{
			ID:        author.ID,
			Username:  author.Username,
			AvatarURL: avatarURL,
		},
	}, nil
}

func (s *tweetService) ReplyToTweetWithMedia(ctx context.Context, userID, tweetID int, tweet *dto.CreateTweetRequest, file *dto.FileData) (*dto.TweetResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	domainTweet := entity.Tweet{
		UserID:        userID,
		ParentTweetID: tweetID,
		Content:       tweet.Content,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	createdTweet, err := s.storage.CreateTweetTx(ctx, tx, &domainTweet)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("user creation failed: %w", err)
	}

	mediaURL, err := s.media.UploadAndAttachTweetMediaTx(ctx, createdTweet.ID, file.File, file.Header.Filename, tx)

	if err := tx.Commit(); err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("commit transaction failed: %w", err)
	}

	author, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("failed to get tweet author")
	}

	avatarURL, err := s.media.GetAvatarUrlByUserID(ctx, userID)
	if err != nil {
		return &dto.TweetResponse{}, fmt.Errorf("failed to get user avatar")
	}

	return &dto.TweetResponse{
		ID:            createdTweet.ID,
		Content:       createdTweet.Content,
		CreatedAt:     createdTweet.CreatedAt,
		UpdatedAt:     createdTweet.UpdatedAt,
		ParentTweetID: &createdTweet.ParentTweetID,
		MediaURL:      mediaURL,
		Author: dto.SmallUserResponse{
			ID:        author.ID,
			Username:  author.Username,
			AvatarURL: avatarURL,
		},
	}, nil
}

func (s *tweetService) GetRepliesToTweet(ctx context.Context, tweetID int) ([]dto.TweetResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	replies, err := s.storage.GetRepliesToParentTweet(ctx, tweetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get replies: %w", err)
	}

	res := make([]dto.TweetResponse, 0, len(replies))
	for _, r := range replies {
		tr, err := s.TweetResponseByTweet(ctx, &r)
		if err != nil {
			return nil, fmt.Errorf("failed to change tweet entity to response: %w", err)
		}
		res = append(res, *tr)
	}
	return res, nil
}
