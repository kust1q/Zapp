package tweets

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

func (s *tweetService) GetTweetById(ctx context.Context, tweetID int) (*entity.Tweet, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tweet, err := s.db.GetTweetById(ctx, tweetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTweetNotFound
		}
		return nil, fmt.Errorf("failed to get tweet by id: %w", err)
	}

	return s.buildEntityTweetToResponse(ctx, tweet)
}

func (s *tweetService) GetTweetsAndRetweetsByUsername(ctx context.Context, username string) ([]entity.Tweet, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tweets, err := s.db.GetTweetsAndRetweetsByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []entity.Tweet{}, nil
		}
		return nil, fmt.Errorf("failed to get tweets by username: %w", err)
	}

	//res := make([]entity.Tweet, 0, len(tweets))
	for i := range tweets {
		processedTweet, err := s.buildEntityTweetToResponse(ctx, &tweets[i])
		if err != nil {
			return nil, err
		}
		tweets[i] = *processedTweet
	}
	return tweets, nil
}

func (s *tweetService) GetRepliesToTweet(ctx context.Context, tweetID int) ([]entity.Tweet, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	replies, err := s.db.GetRepliesToTweet(ctx, tweetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get replies: %w", err)
	}

	for i := range replies {
		processedTweet, err := s.buildEntityTweetToResponse(ctx, &replies[i])
		if err != nil {
			return nil, err
		}
		replies[i] = *processedTweet
	}
	return replies, nil
}
