package tweets

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

func (s *tweetService) UpdateTweet(ctx context.Context, req *entity.Tweet) (*entity.Tweet, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	tweet, err := s.db.GetTweetById(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTweetNotFound
		}
		return nil, fmt.Errorf("failed to get tweet by id: %w", err)
	}
	if tweet.Author.ID != req.Author.ID {
		return nil, ErrUnauthorizedUpdate
	}
	tweet.Content = req.Content
	tweet.UpdatedAt = time.Now()
	updatedTweet, err := s.storage.UpdateTweet(ctx, tweet)
	if err != nil {
		return nil, fmt.Errorf("failed to update tweet: %w", err)
	}

	return &dto.UpdateTweetResponse{
		ID:        updatedTweet.ID,
		UserID:    updatedTweet.UserID,
		Content:   updatedTweet.Content,
		UpdatedAt: updatedTweet.UpdatedAt,
	}, nil
}
