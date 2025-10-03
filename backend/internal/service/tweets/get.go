package tweets

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/dto"
)

func (s *tweetService) GetTweetById(ctx context.Context, tweetID int) (*dto.TweetResponseWithCounters, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tweet, err := s.storage.GetTweetById(ctx, tweetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &dto.TweetResponseWithCounters{}, ErrTweetNotFound
		}
		return &dto.TweetResponseWithCounters{}, fmt.Errorf("failed to get tweet by id: %w", err)
	}
	return s.tweetResponseWithCountersByTweet(ctx, tweet)
}

func (s *tweetService) GetTweetsAndRetweetsByUsername(ctx context.Context, username string) ([]dto.TweetResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tweets, err := s.storage.GetTweetsAndRetweetsByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []dto.TweetResponse{}, nil
		}
		return nil, fmt.Errorf("failed to get tweets by username: %w", err)
	}

	res := make([]dto.TweetResponse, 0, len(tweets))
	for _, t := range tweets {
		tr, err := s.TweetResponseByTweet(ctx, &t)
		if err != nil {
			return nil, fmt.Errorf("failed to get tweet responses by username: %w", err)
		}
		res = append(res, *tr)
	}
	return res, nil
}
