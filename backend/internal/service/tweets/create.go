package tweets

import (
	"context"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/sirupsen/logrus"
)

func (s *tweetService) CreateTweet(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	createdTweet, err := s.db.CreateTweetTx(ctx, tx, tweet)
	if err != nil {
		return nil, fmt.Errorf("user creation failed: %w", err)
	}

	var mediaUrl string
	if tweet.File != nil {
		mediaUrl, err = s.media.UploadAndAttachTweetMediaTx(ctx, createdTweet.ID, tweet.File.File, tweet.File.Header.Filename, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to upload media: %w", err)
		}
	}

	createdTweet.MediaUrl = mediaUrl

	response, err := s.BuildEntityTweetToResponse(ctx, createdTweet)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction failed: %w", err)
	}

	go func(t *entity.Tweet) {
		cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.search.IndexTweet(cntx, t); err != nil {
			logrus.WithError(err).
				WithField("tweet_id", t.ID).
				Warn("failed to index tweet in elastic")
		}
	}(createdTweet)

	return response, nil
}
