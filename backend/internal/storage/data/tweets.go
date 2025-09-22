package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/dto"
	"github.com/kust1q/Zapp/backend/internal/storage/postgres"
)

func (s *dataStorage) CreateTweet(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, parent_tweet_id, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", postgres.TweetsTable)
	var id int
	if err := s.db.QueryRowContext(ctx, query, tweet.UserID, tweet.ParentTweetID, tweet.Content, tweet.CreatedAt, tweet.UpdatedAt).Scan(&id); err != nil {
		return &entity.Tweet{}, err
	}
	tweet.ID = id
	return tweet, nil
}

func (s *dataStorage) CreateTweetTx(ctx context.Context, tx *sql.Tx, tweet *entity.Tweet) (*entity.Tweet, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, parent_tweet_id, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", postgres.TweetsTable)
	var id int
	if err := tx.QueryRowContext(ctx, query, tweet.UserID, tweet.ParentTweetID, tweet.Content, tweet.CreatedAt, tweet.UpdatedAt).Scan(&id); err != nil {
		return &entity.Tweet{}, err
	}
	tweet.ID = id
	return tweet, nil
}

func (s *dataStorage) GetTweetById(ctx context.Context, tweetID int) (*entity.Tweet, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND user_id = $2", postgres.TweetsTable)
	var tweet entity.Tweet
	err := s.db.GetContext(ctx, &tweet, query, tweetID)
	return &tweet, err
}

func (s *dataStorage) UpdateTweet(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error) {
	query := fmt.Sprintf("UPDATE %s SET content = $1, updated_at = $2 WHERE id = $3 RETURNING content, updated_at", postgres.TweetsTable)

	var content string
	var updatedAt time.Time
	err := s.db.QueryRowContext(ctx, query, tweet.Content, tweet.UpdatedAt, tweet.ID).Scan(&content, &updatedAt)
	if err != nil {
		return &entity.Tweet{}, err
	}

	tweet.Content = content
	tweet.UpdatedAt = updatedAt
	return tweet, nil
}

func (s *dataStorage) DeleteTweet(ctx context.Context, userID, tweetID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", postgres.TweetsTable)
	result, err := s.db.ExecContext(ctx, query, tweetID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("tweet not found")
	}
	return nil
}

func (s *dataStorage) LikeTweet(ctx context.Context, userID, tweetID int) error {
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

func (s *dataStorage) UnLikeTweet(ctx context.Context, userID, tweetID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND tweet_id = $2", postgres.LikesTable)
	result, err := s.db.ExecContext(ctx, query, userID, tweetID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("tweet not found")
	}
	return nil
}

func (s *dataStorage) Retweet(ctx context.Context, userID, tweetID int, createdAt time.Time) error {
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

func (s *dataStorage) DeleteRetweet(ctx context.Context, userID, retweetID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND tweet_id = $2", postgres.RetweetsTable)
	result, err := s.db.ExecContext(ctx, query, userID, retweetID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("retweet not found")
	}

	return nil
}

func (s *dataStorage) GetRepliesToParentTweet(ctx context.Context, parentTweetID int) ([]entity.Tweet, error) {
	query := fmt.Sprintf("SELECT id, user_id, parent_tweet_id, content, created_at, updated_at FROM %s WHERE parent_tweet_id = $1 ORDER BY created_at DESC", postgres.TweetsTable)
	var tweets []entity.Tweet
	if err := s.db.SelectContext(ctx, &tweets, query, parentTweetID); err != nil {
		return nil, err
	}
	return tweets, nil
}

func (s *dataStorage) GetTweetsAndRetweetsByUsername(ctx context.Context, username string) ([]entity.Tweet, error) {
	query := fmt.Sprintf(`
        SELECT t.id, t.user_id, t.parent_tweet_id, t.content, t.created_at, t.updated_at 
        FROM %s t 
        JOIN %s u ON t.user_id = u.id 
        WHERE u.username = $1
		UNION ALL
		SELECT t.id, t.user_id, t.parent_tweet_id, t.content, r.created_at, t.updated_at
		FROM %s r
		JOIN %s t ON r.tweet_id = t.id 
		JOIN %s u ON r.user_id = u.id
		WHERE u.username = $1 
		ORDER BY r.created_at DESC
		`,
		postgres.TweetsTable, postgres.UserTable, postgres.RetweetsTable, postgres.TweetsTable, postgres.UserTable)

	var tweets []entity.Tweet
	if err := s.db.SelectContext(ctx, &tweets, query, username); err != nil {
		return nil, fmt.Errorf("failed to get tweets by username: %w", err)
	}
	return tweets, nil
}

func (s *dataStorage) GetCounts(ctx context.Context, tweetID int) (likes, retweets, replies int, err error) {
	query := fmt.Sprintf(`
        SELECT 
            (SELECT COUNT(*) FROM %s WHERE tweet_id = $1) as likes,
            (SELECT COUNT(*) FROM %s WHERE tweet_id = $1) as retweets,
            (SELECT COUNT(*) FROM %s WHERE parent_tweet_id = $1) as replies
    `, postgres.LikesTable, postgres.RetweetsTable, postgres.TweetsTable)

	err = s.db.QueryRowContext(ctx, query, tweetID).Scan(&likes, &retweets, &replies)
	return
}

func (s *dataStorage) GetLikes(ctx context.Context, tweetID int) ([]dto.UserLikeResponse, error) {
	query := fmt.Sprintf(`
        SELECT u.id as user_id, u.username, u.avatar_url
        FROM %s l
        JOIN %s u ON l.user_id = u.id
        WHERE l.tweet_id = $1`,
		postgres.LikesTable, postgres.UserTable)

	var userLikes []dto.UserLikeResponse
	if err := s.db.SelectContext(ctx, &userLikes, query, tweetID); err != nil {
		return nil, fmt.Errorf("failed to get users who liked tweet: %w", err)
	}
	return userLikes, nil
}
