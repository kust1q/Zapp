package user

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *userService) DeleteUser(ctx context.Context, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := s.media.DeleteAvatar(ctx, userID); err != nil {
		logrus.WithField("user_id", userID).Warnf("failed to delete user avatar: %s", err)
	}
	return s.storage.DeleteUser(ctx, userID)
}
