package tweets

import (
	"context"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/sirupsen/logrus"
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

func (s *tweetService) GetLikes(ctx context.Context, tweetID int) ([]entity.SmallUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	users, err := s.db.GetLikes(ctx, tweetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get likes: %w", err)
	}

	for i := range users {
		users[i].AvatarURL, err = s.media.GetAvatarUrlByUserID(ctx, users[i].ID)
		if err != nil {
			logrus.WithError(err).Warn("failed to get avatar url")
		}
	}
	return users, nil
}
