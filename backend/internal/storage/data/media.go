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
		INSERT INTO %s (tweet_id, path, mime_type, size_bytes)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (tweet_id) DO UPDATE
		SET path = EXCLUDED.path,
			mime_type = EXCLUDED.mime_type,
			size_bytes = EXCLUDED.size_bytes
		RETURNING id
	`, postgres.TweetMediaTable)
	var id int
	if err := tx.QueryRowContext(ctx, query, media.TweetID, media.Path, media.MimeType, media.SizeBytes).Scan(&id); err != nil {
		return &entity.TweetMedia{}, err
	}
	media.ID = id
	return media, nil
}

func (s *dataStorage) GetMediaUrlByTweetID(ctx context.Context, tweetID int) (string, error) {
	query := fmt.Sprintf("SELECT path FROM %s WHERE tweet_id = $1", postgres.TweetMediaTable)
	var path string
	if err := s.db.GetContext(ctx, &path, query, tweetID); err != nil {
		return "", err
	}
	return path, nil
}

func (s *dataStorage) GetMediaDataByTweetID(ctx context.Context, tweetID int) (*entity.TweetMedia, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE tweet_id = $1", postgres.TweetMediaTable)
	var media entity.TweetMedia
	if err := s.db.GetContext(ctx, &media, query, tweetID); err != nil {
		return &entity.TweetMedia{}, err
	}
	return &media, nil
}

func (s *dataStorage) DeleteMediaByTweetID(ctx context.Context, tweetID, userID int) error {
	query := fmt.Sprintf(`
        DELETE FROM %s m
        WHERE m.tweet_id = $1 
        AND EXISTS (
            SELECT 1 FROM %s t 
            WHERE t.id = m.tweet_id AND t.user_id = $2)`,
		postgres.TweetMediaTable,
		postgres.TweetsTable)

	result, err := s.db.ExecContext(ctx, query, tweetID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("tweet media not found")
	}

	return nil
}

func (s *dataStorage) UploadAvatarTx(ctx context.Context, tx *sql.Tx, avatar *entity.Avatar) (*entity.Avatar, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, path, mime_type, size_bytes) VALUES ($1, $2, $3, $4) ON CONFLICT (user_id) DO NOTHING RETURNING id", postgres.AvatarsTable)
	var id int
	if err := tx.QueryRowContext(ctx, query, avatar.UserID, avatar.Path, avatar.MimeType, avatar.SizeBytes).Scan(&id); err != nil {
		return &entity.Avatar{}, err
	}
	avatar.ID = id
	return avatar, nil
}

func (s *dataStorage) GetAvatarUrlByUserID(ctx context.Context, userID int) (string, error) {
	query := fmt.Sprintf("SELECT path FROM %s WHERE user_id = $1", postgres.AvatarsTable)
	var path string
	if err := s.db.GetContext(ctx, path, query, userID); err != nil {
		return "", err
	}
	return path, nil
}

func (s *dataStorage) GetAvatarDataByUserID(ctx context.Context, userID int) (*entity.Avatar, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", postgres.AvatarsTable)
	var avatar entity.Avatar
	if err := s.db.GetContext(ctx, &avatar, query, userID); err != nil {
		return &entity.Avatar{}, err
	}
	return &avatar, nil
}

func (s *dataStorage) DeleteAvatarByUserID(ctx context.Context, userID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", postgres.AvatarsTable)
	result, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("avatar not found")
	}

	return nil
}
