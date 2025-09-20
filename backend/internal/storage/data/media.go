package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/storage/postgres"
)

func (s *dataStorage) UpsertByTweetIdTx(ctx context.Context, tx *sql.Tx, media *entity.TweetMedia) (*entity.TweetMedia, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (tweet_id, media_url, mime_type, size_bytes)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (tweet_id) DO UPDATE
		SET media_url = EXCLUDED.media_url,
			mime_type = EXCLUDED.mime_type,
			size_bytes = EXCLUDED.size_bytes
		RETURNING id
	`, postgres.TweetMediaTable)
	var id int
	if err := tx.QueryRowContext(ctx, query, media.TweetID, media.MediaURL, media.MimeType, media.SizeBytes).Scan(&id); err != nil {
		return &entity.TweetMedia{}, err
	}
	media.ID = id
	return media, nil
}

func (s *dataStorage) GetMediaByTweetID(ctx context.Context, tweetID int) (*entity.TweetMedia, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE tweet_id = $1", postgres.TweetMediaTable)
	var media entity.TweetMedia
	if err := s.db.GetContext(ctx, &media, query, tweetID); err != nil {
		return &entity.TweetMedia{}, err
	}
	return &media, nil
}

func (s *dataStorage) DeleteMediaByTweetID(ctx context.Context, tweetID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE tweet_id = $1", postgres.TweetMediaTable)
	_, err := s.db.ExecContext(ctx, query, tweetID)
	return err
}

func (s *dataStorage) UploadAvatarTx(ctx context.Context, tx *sql.Tx, avatar *entity.Avatar) (*entity.Avatar, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, media_url, mime_type, size_bytes) VALUES ($1, $2, $3, $4) ON CONFLICT (user_id) DO NOTHING RETURNING id", postgres.AvatarsTable)
	var id int
	if err := tx.QueryRowContext(ctx, query, avatar.UserID, avatar.MediaURL, avatar.MimeType, avatar.SizeBytes).Scan(&id); err != nil {
		return &entity.Avatar{}, err
	}
	avatar.ID = id
	return avatar, nil
}

func (s *dataStorage) GetAvatarByUserID(ctx context.Context, userID int) (*entity.Avatar, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", postgres.AvatarsTable)
	var avatar entity.Avatar
	if err := s.db.GetContext(ctx, &avatar, query, userID); err != nil {
		return &entity.Avatar{}, err
	}
	return &avatar, nil
}
