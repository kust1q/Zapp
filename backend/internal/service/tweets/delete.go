package tweets

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *tweetService) DeleteTweet(ctx context.Context, userID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tweet, err := s.storage.GetTweetById(ctx, tweetID)
	if err != nil {
		return fmt.Errorf("failed to get tweet by id: %w", err)
	}
	if tweet.UserID != userID {
		return ErrNotEnoughRights
	}
	if err := s.media.DeleteTweetMedia(ctx, tweetID); err != nil {
		logrus.WithField("tweet_id", tweetID).Warnf("failed to delete tweet media")
	}
	return s.storage.DeleteTweet(ctx, userID, tweetID)
}
