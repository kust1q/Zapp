package tweets

import (
	"context"
	"time"

	"github.com/kust1q/Zapp/backend/internal/dto"
)

func (s *tweetService) LikeTweet(ctx context.Context, userID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.storage.LikeTweet(ctx, userID, tweetID)
}

func (s *tweetService) UnlikeTweet(ctx context.Context, userID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.storage.UnLikeTweet(ctx, userID, tweetID)
}

func (s *tweetService) GetLikes(ctx context.Context, tweetID int) ([]dto.UserLikeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.storage.GetLikes(ctx, tweetID)
}
