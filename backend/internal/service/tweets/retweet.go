package tweets

import (
	"context"
	"time"
)

func (s *tweetService) CreateRetweet(ctx context.Context, userID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.storage.Retweet(ctx, userID, tweetID, time.Now())
}

func (s *tweetService) DeleteRetweet(ctx context.Context, userID, retweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.storage.DeleteRetweet(ctx, userID, retweetID)
}
