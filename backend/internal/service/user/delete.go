package user

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *userService) DeleteUser(ctx context.Context, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := s.media.DeleteAvatar(ctx, userID); err != nil {
		logrus.WithField("user_id", userID).Warnf("failed to delete user avatar: %s", err)
	}
	if err := s.db.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user")
	}
	go func(id int) {
		cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.search.DeleteUser(cntx, id); err != nil {
			logrus.WithError(err).WithField("user_id", id).Warn("failed to delete user from elastic")
		}
		if err := s.search.DeleteTweetsByUserID(cntx, id); err != nil {
			logrus.WithError(err).WithField("user_id", id).Warn("failed to delete user tweets from elastic")
		}
	}(userID)
	return nil
}
