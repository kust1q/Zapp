package tweets

import (
	"context"
	"time"
)

func (s *tweetService) DeleteTweet(ctx context.Context, userID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err := s.storage.DeleteTweet(ctx, userID, tweetID)
	return err
}
