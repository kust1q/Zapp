package tweets

import (
	"context"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

func (s *tweetService) LikeTweet(ctx context.Context, userID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.db.LikeTweet(ctx, userID, tweetID)
}

func (s *tweetService) UnlikeTweet(ctx context.Context, userID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.db.UnLikeTweet(ctx, userID, tweetID)
}

func (s *tweetService) GetLikes(ctx context.Context, tweetID int) ([]entity.Like, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.db.GetLikes(ctx, tweetID)
}
