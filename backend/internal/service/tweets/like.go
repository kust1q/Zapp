package tweets

import (
	"context"
	"time"
)

func (s *tweetService) LikeTweet(ctx context.Context, userID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.storage.LikeTweet(ctx, userID, tweetID)
}

func (s *tweetService) UnLikeTweet(ctx context.Context, userID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.storage.UnLikeTweet(ctx, userID, tweetID)
}
