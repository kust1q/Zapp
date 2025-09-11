package tweets

import (
	"context"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/dto"
)

func (s *tweetService) ReplyToTweet(ctx context.Context, userID, tweetID int, tweet *dto.CreateTweetRequest) (dto.TweetResponse, error) {
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
		return dto.TweetResponse{}, fmt.Errorf("tweet creation failed: %w", err)
	}

	return dto.TweetResponse{
		ID:            createdTweet.ID,
		UserID:        createdTweet.UserID,
		Content:       createdTweet.Content,
		CreatedAt:     createdTweet.CreatedAt,
		UpdatedAt:     createdTweet.UpdatedAt,
		ParentTweetID: &createdTweet.ParentTweetID,
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
