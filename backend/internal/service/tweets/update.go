package tweets

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/sirupsen/logrus"
)

func (s *tweetService) UpdateTweet(ctx context.Context, req *entity.Tweet) (*entity.Tweet, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	exTweet, err := s.db.GetTweetById(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTweetNotFound
		}
		return nil, fmt.Errorf("failed to get tweet by id: %w", err)
	}

	if exTweet.Author.ID != req.Author.ID {
		return nil, ErrUnauthorizedUpdate
	}

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	exTweet.Content = req.Content
	exTweet.UpdatedAt = time.Now()

	updatedTweet, err := s.db.UpdateTweet(ctx, exTweet)
	if err != nil {
		return nil, fmt.Errorf("failed to update tweet: %w", err)
	}

	var mediaUrl string
	if req.File != nil {
		if err := s.media.DeleteTweetMedia(ctx, req.ID, req.Author.ID); err != nil {
			logrus.WithError(err).WithField("tweet_id", req.ID).Warn("failed to delete old media record")
		}
		mediaUrl, err = s.media.UploadAndAttachTweetMediaTx(ctx, updatedTweet.ID, req.File.File, req.File.Header.Filename, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to upload new media: %w", err)
		}
	}

	updatedTweet.MediaUrl = mediaUrl

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	go func(t *entity.Tweet) {
		cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.search.IndexTweet(cntx, t); err != nil {
			logrus.WithError(err).
				WithField("tweet_id", t.ID).
				Warn("failed to update tweet index in elastic")
		}
	}(updatedTweet)

	return s.BuildEntityTweetToResponse(ctx, updatedTweet)
}
