package tweets

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/dto"
)

var (
	ErrTweetNotFound      = errors.New("tweet not found")
	ErrUnauthorizedUpdate = errors.New("user is not authorized to update this tweet")
)

func (s *tweetService) UpdateTweet(ctx context.Context, userID, tweetID int, req *dto.UpdateTweetRequest) (dto.UpdateTweetResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	tweet, err := s.storage.GetTweetById(ctx, tweetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.UpdateTweetResponse{}, ErrTweetNotFound
		}
		return dto.UpdateTweetResponse{}, fmt.Errorf("failed to get tweet by id: %w", err)
	}
	if tweet.UserID != userID {
		return dto.UpdateTweetResponse{}, ErrUnauthorizedUpdate
	}
	tweet.Content = req.Content
	tweet.UpdatedAt = time.Now()
	updatedTweet, err := s.storage.UpdateTweet(ctx, &tweet)
	if err != nil {
		return dto.UpdateTweetResponse{}, fmt.Errorf("failed to update tweet: %w", err)
	}

	return dto.UpdateTweetResponse{
		ID:        updatedTweet.ID,
		UserID:    updatedTweet.UserID,
		Content:   updatedTweet.Content,
		UpdatedAt: updatedTweet.UpdatedAt,
	}, nil
}
