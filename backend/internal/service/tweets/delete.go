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
	if err := s.media.DeleteTweetMedia(ctx, tweetID, userID); err != nil {
		logrus.WithField("tweet_id", tweetID).Warnf("failed to delete tweet media")
	}
	if err := s.db.DeleteTweet(ctx, userID, tweetID); err != nil {
		return fmt.Errorf("failed to delete tweet")
	}
	go func(id int) {
		cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.search.DeleteTweet(cntx, id); err != nil {
			logrus.WithError(err).WithField("tweet_id", id).Warn("failed to delete tweet from elastic")
		}
	}(tweetID)
	return nil
}
