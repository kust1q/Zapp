package tweets

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *tweetService) DeleteTweet(ctx context.Context, userID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := s.media.DeleteTweetMedia(ctx, tweetID, userID); err != nil {
		logrus.WithField("tweet_id", tweetID).Warnf("failed to delete tweet media")
	}
	return s.storage.DeleteTweet(ctx, userID, tweetID)
}
