package data

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/storage/postgres"
)

type tweetStorage struct {
	db *sqlx.DB
}

func NewTweetStorage(db *sqlx.DB) *tweetStorage {
	return &tweetStorage{
		db: db,
	}
}

func (s *tweetStorage) CreateTweet(ctx context.Context, tweet *entity.Tweet) (entity.Tweet, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, parent_tweet_id, content, created_at, updated_at) VALUES (:user_id, :parent_tweet_id, :content, :created_at, :updated_at) RETURNING id", postgres.TweetsTable)

	stmt, err := s.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return entity.Tweet{}, err
	}

	err = stmt.QueryRowContext(ctx, tweet).Scan(&tweet.ID)
	if err != nil {
		return entity.Tweet{}, err
	}

	return *tweet, nil
}

func (s *tweetStorage) GetTweetByIds(ctx context.Context, tweetID, userID int) (entity.Tweet, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND user_id = $2", postgres.TweetsTable)
	var tweet entity.Tweet
	err := s.db.GetContext(ctx, &tweet, query, tweetID, userID)
	return tweet, err
}

func (s *tweetStorage) UpdateTweet(ctx context.Context, userID int, tweet *entity.Tweet) (entity.Tweet, error) {
	query := fmt.Sprintf("UPDATE %s SET content = $1, updated_at = $2 WHERE id = $3 AND user_id = $4 RETURNING content, updated_at", postgres.TweetsTable)

	var content string
	var updatedAt time.Time
	err := s.db.QueryRowContext(ctx, query, tweet.Content, tweet.UpdatedAt, tweet.ID, userID).Scan(&content, &updatedAt)
	if err != nil {
		return entity.Tweet{}, err
	}

	tweet.Content = content
	tweet.UpdatedAt = updatedAt
	return *tweet, nil
}

func (s *tweetStorage) DeleteTweet(ctx context.Context, userID, tweetID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", postgres.TweetsTable)
	result, err := s.db.ExecContext(ctx, query, tweetID, userID)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("tweet not found")
	}

	return nil
}

func (s *tweetStorage) LikeTweet(ctx context.Context, userID, tweetID int) error {
	var exists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", postgres.TweetsTable)
	err := s.db.QueryRowContext(ctx, checkQuery, tweetID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("tweet not found")
	}

	query := fmt.Sprintf("INSERT INTO %s (user_id, tweet_id) VALUES ($1, $2) ON CONFLICT (user_id, tweet_id) DO NOTHING", postgres.LikesTable)
	_, err = s.db.ExecContext(ctx, query, userID, tweetID)
	return err
}

func (s *tweetStorage) UnLikeTweet(ctx context.Context, userID, tweetID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND tweet_id = $2", postgres.LikesTable)
	result, err := s.db.ExecContext(ctx, query, userID, tweetID)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("tweet not found")
	}

	return nil
}

func (s *tweetStorage) Retweet(ctx context.Context, userID, tweetID int, createdAt time.Time) error {
	var exists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", postgres.TweetsTable)
	err := s.db.QueryRowContext(ctx, checkQuery, tweetID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("tweet not found")
	}

	query := fmt.Sprintf("INSERT INTO %s (user_id, tweet_id, created_at) VALUES ($1, $2, $3)", postgres.RetweetsTable)
	_, err = s.db.ExecContext(ctx, query, userID, tweetID, createdAt)
	return err
}

func (s *tweetStorage) DeleteRetweet(ctx context.Context, userID, retweetID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND tweet_id = $2", postgres.RetweetsTable)
	result, err := s.db.ExecContext(ctx, query, userID, retweetID)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("retweet not found")
	}

	return nil
}

func (s *tweetStorage) GetLikeCount(ctx context.Context, tweetID int) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE tweet_id = $1", postgres.LikesTable)
	var count int
	err := s.db.QueryRowContext(ctx, query, tweetID).Scan(&count)
	return count, err
}

func (s *tweetStorage) GetRetweetCount(ctx context.Context, tweetID int) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE tweet_id = $1", postgres.RetweetsTable)
	var count int
	err := s.db.QueryRowContext(ctx, query, tweetID).Scan(&count)
	return count, err
}

func (s *tweetStorage) GetReplyCount(ctx context.Context, tweetID int) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE parent_tweet_id = $1", postgres.TweetsTable)
	var count int
	err := s.db.QueryRowContext(ctx, query, tweetID).Scan(&count)
	return count, err
}

func (s *tweetStorage) GetCounts(ctx context.Context, tweetID int) (likes, retweets, replies int, err error) {
	query := fmt.Sprintf(`
        SELECT 
            (SELECT COUNT(*) FROM %s WHERE tweet_id = $1) as likes,
            (SELECT COUNT(*) FROM %s WHERE tweet_id = $1) as retweets,
            (SELECT COUNT(*) FROM %s WHERE parent_tweet_id = $1) as replies
    `, postgres.LikesTable, postgres.RetweetsTable, postgres.TweetsTable)

	err = s.db.QueryRowContext(ctx, query, tweetID).Scan(&likes, &retweets, &replies)
	return
}
