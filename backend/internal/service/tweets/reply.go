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
		UserID:        userID,
		Content:       createdTweet.Content,
		CreatedAt:     createdTweet.CreatedAt,
		UpdatedAt:     createdTweet.UpdatedAt,
		ParentTweetID: &createdTweet.ParentTweetID,
	}, nil
}
