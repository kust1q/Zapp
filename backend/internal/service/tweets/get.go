package tweets

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/errs"
)

func (s *tweetService) GetTweetById(ctx context.Context, tweetID int) (*entity.Tweet, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tweet, err := s.db.GetTweetById(ctx, tweetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrTweetNotFound
		}
		return nil, fmt.Errorf("failed to get tweet by id: %w", err)
	}
	if tweet.MediaUrl != "" {
		mediaUrl, err := s.media.GetMediaUrlByTweetID(ctx, tweet.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get tweet media url: %w", err)
		}

		tweet.MediaUrl = mediaUrl
	}

	return s.BuildEntityTweetToResponse(ctx, tweet)
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

	for i := range tweets {
		if tweets[i].MediaUrl != "" {
			mediaUrl, err := s.media.GetMediaUrlByTweetID(ctx, tweets[i].ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get tweet media url")
			}
			tweets[i].MediaUrl = mediaUrl
		}
		processedTweet, err := s.BuildEntityTweetToResponse(ctx, &tweets[i])
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
		if replies[i].MediaUrl != "" {
			mediaUrl, err := s.media.GetMediaUrlByTweetID(ctx, replies[i].ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get tweet media url")
			}
			replies[i].MediaUrl = mediaUrl
		}

		processedTweet, err := s.BuildEntityTweetToResponse(ctx, &replies[i])
		if err != nil {
			return nil, err
		}
		replies[i] = *processedTweet
	}
	return replies, nil
}
