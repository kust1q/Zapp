package tweets

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/dto"
)

func (s *tweetService) GetTweetById(ctx context.Context, tweetID int) (dto.TweetResponseWithCounters, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tweet, err := s.storage.GetTweetById(ctx, tweetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.TweetResponseWithCounters{}, ErrTweetNotFound
		}
		return dto.TweetResponseWithCounters{}, fmt.Errorf("failed to get tweet by id: %w", err)
	}

	likes, retweets, replyCount, err := s.storage.GetCounts(ctx, tweetID)
	return dto.TweetResponseWithCounters{
		TweetResponse: dto.TweetResponse{
			ID:            tweet.ID,
			UserID:        tweet.UserID,
			Content:       tweet.Content,
			CreatedAt:     tweet.CreatedAt,
			UpdatedAt:     tweet.UpdatedAt,
			ParentTweetID: &tweet.ParentTweetID,
		},
		ReplyCount:   replyCount,
		RetweetCount: retweets,
		LikeCount:    likes,
	}, nil
}

func (s *tweetService) GetTweetsByUsername(ctx context.Context, username string) ([]dto.TweetResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tweets, err := s.storage.GetTweetsByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []dto.TweetResponse{}, nil
		}
		return []dto.TweetResponse{}, fmt.Errorf("failed to get tweets by username: %w", err)
	}

	res := make([]dto.TweetResponse, 0, len(tweets))
	for _, r := range tweets {
		res = append(res, dto.TweetResponse{
			ID:            r.ID,
			UserID:        r.UserID,
			Content:       r.Content,
			CreatedAt:     r.CreatedAt,
			UpdatedAt:     r.UpdatedAt,
			ParentTweetID: &r.ParentTweetID,
		})
	}
	return res, nil
}
