package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	conv "github.com/kust1q/Zapp/backend/internal/pkg/conv/db"
	"github.com/kust1q/Zapp/backend/internal/providers/db/models"
	"github.com/kust1q/Zapp/backend/internal/storage/postgres"
)

func (pg *PostgresDB) CreateTweet(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error) {
	tweetModel := conv.FromDomainToTweetModel(tweet)
	if tweetModel == nil {
		return nil, fmt.Errorf("cannot convert nil entity to DB model")
	}

	query := fmt.Sprintf("INSERT INTO %s (user_id, parent_tweet_id, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", postgres.TweetsTable)
	var id int
	if err := pg.db.QueryRowContext(ctx, query, tweetModel.UserID, tweetModel.ParentTweetID, tweetModel.Content, tweetModel.CreatedAt, tweetModel.UpdatedAt).Scan(&id); err != nil {
		return nil, err
	}
	tweetModel.ID = id
	createdtweet := conv.FromTweetModelToDomain(tweetModel)
	return createdtweet, nil
}

func (pg *PostgresDB) CreateTweetTx(ctx context.Context, tx *sql.Tx, tweet *entity.Tweet) (*entity.Tweet, error) {
	tweetModel := conv.FromDomainToTweetModel(tweet)
	if tweetModel == nil {
		return nil, fmt.Errorf("cannot convert nil entity to DB model")
	}

	query := fmt.Sprintf("INSERT INTO %s (user_id, parent_tweet_id, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", postgres.TweetsTable)
	var id int
	if err := tx.QueryRowContext(ctx, query, tweetModel.UserID, tweetModel.ParentTweetID, tweetModel.Content, tweetModel.CreatedAt, tweetModel.UpdatedAt).Scan(&id); err != nil {
		return nil, err
	}
	tweetModel.ID = id
	createdtweet := conv.FromTweetModelToDomain(tweetModel)
	return createdtweet, nil
}

func (pg *PostgresDB) GetTweetById(ctx context.Context, tweetID int) (*entity.Tweet, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", postgres.TweetsTable)
	var tweetModel models.Tweet
	err := pg.db.GetContext(ctx, &tweetModel, query, tweetID)
	if err != nil {
		return nil, err
	}
	tweet := conv.FromTweetModelToDomain(&tweetModel)
	return tweet, nil
}

func (pg *PostgresDB) UpdateTweet(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error) {
	query := fmt.Sprintf("UPDATE %s SET content = $1, updated_at = $2 WHERE id = $3 RETURNING content, updated_at", postgres.TweetsTable)
	var content string
	var updatedAt time.Time
	err := pg.db.QueryRowContext(ctx, query, tweet.Content, tweet.UpdatedAt, tweet.ID).Scan(&content, &updatedAt)
	if err != nil {
		return nil, err
	}
	tweet.Content = content
	tweet.UpdatedAt = updatedAt
	return tweet, nil
}

func (pg *PostgresDB) DeleteTweet(ctx context.Context, userID, tweetID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", postgres.TweetsTable)
	result, err := pg.db.ExecContext(ctx, query, tweetID, userID)
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

func (pg *PostgresDB) LikeTweet(ctx context.Context, userID, tweetID int) error {
	var exists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", postgres.TweetsTable)
	err := pg.db.QueryRowContext(ctx, checkQuery, tweetID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("tweet not found")
	}
	query := fmt.Sprintf("INSERT INTO %s (user_id, tweet_id) VALUES ($1, $2) ON CONFLICT (user_id, tweet_id) DO NOTHING", postgres.LikesTable)
	_, err = pg.db.ExecContext(ctx, query, userID, tweetID)
	return err
}

func (pg *PostgresDB) UnLikeTweet(ctx context.Context, userID, tweetID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND tweet_id = $2", postgres.LikesTable)
	result, err := pg.db.ExecContext(ctx, query, userID, tweetID)
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

func (pg *PostgresDB) Retweet(ctx context.Context, userID, tweetID int, createdAt time.Time) error {
	var exists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", postgres.TweetsTable)
	err := pg.db.QueryRowContext(ctx, checkQuery, tweetID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("tweet not found")
	}
	query := fmt.Sprintf("INSERT INTO %s (user_id, tweet_id, created_at) VALUES ($1, $2, $3)", postgres.RetweetsTable)
	_, err = pg.db.ExecContext(ctx, query, userID, tweetID, createdAt)
	return err
}

func (pg *PostgresDB) DeleteRetweet(ctx context.Context, userID, retweetID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND tweet_id = $2", postgres.RetweetsTable)
	result, err := pg.db.ExecContext(ctx, query, userID, retweetID)
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

func (pg *PostgresDB) GetRepliesToTweet(ctx context.Context, parentTweetID int) ([]entity.Tweet, error) {
	query := fmt.Sprintf("SELECT id, user_id, parent_tweet_id, content, created_at, updated_at FROM %s WHERE parent_tweet_id = $1 ORDER BY created_at DESC", postgres.TweetsTable)
	var tweetModels []models.Tweet
	if err := pg.db.SelectContext(ctx, &tweetModels, query, parentTweetID); err != nil {
		return nil, err
	}

	return conv.FromTweetModelToDomainList(tweetModels), nil
}

func (pg *PostgresDB) GetTweetsAndRetweetsByUsername(ctx context.Context, username string) ([]entity.Tweet, error) {
	query := fmt.Sprintf(`
        SELECT t.id, t.user_id, t.parent_tweet_id, t.content, t.created_at, t.updated_at 
        FROM %s t 
        JOIN %s u ON t.user_id = u.id 
        WHERE u.username = $1
        UNION ALL
        SELECT t.id, r.user_id, t.parent_tweet_id, t.content, r.created_at, t.updated_at
        FROM %s r
        JOIN %s t ON r.tweet_id = t.id 
        JOIN %s u ON r.user_id = u.id
        WHERE u.username = $1 
        ORDER BY created_at DESC`,
		postgres.TweetsTable, postgres.UserTable, postgres.RetweetsTable, postgres.TweetsTable, postgres.UserTable)

	var tweetModels []models.Tweet
	if err := pg.db.SelectContext(ctx, &tweetModels, query, username); err != nil {
		return nil, fmt.Errorf("failed to get tweets by username: %w", err)
	}

	return conv.FromTweetModelToDomainList(tweetModels), nil
}

func (pg *PostgresDB) GetCounts(ctx context.Context, tweetID int) (likes, retweets, replies int, err error) {
	query := fmt.Sprintf(`
        SELECT 
            (SELECT COUNT(*) FROM %s WHERE tweet_id = $1) as likes,
            (SELECT COUNT(*) FROM %s WHERE tweet_id = $1) as retweets,
            (SELECT COUNT(*) FROM %s WHERE parent_tweet_id = $1) as replies
    `, postgres.LikesTable, postgres.RetweetsTable, postgres.TweetsTable)

	err = pg.db.QueryRowContext(ctx, query, tweetID).Scan(&likes, &retweets, &replies)
	return
}

func (pg *PostgresDB) GetLikes(ctx context.Context, tweetID int) ([]entity.Like, error) {
	query := fmt.Sprintf(`
        SELECT u.id as user_id, u.username, u.avatar_url
        FROM %s l
        JOIN %s u ON l.user_id = u.id
        WHERE l.tweet_id = $1`,
		postgres.LikesTable, postgres.UserTable)

	var userLikes []models.Like
	if err := pg.db.SelectContext(ctx, &userLikes, query, tweetID); err != nil {
		return nil, fmt.Errorf("failed to get users who liked tweet: %w", err)
	}
	return conv.FromLikeModelToDomainList(userLikes), nil
}
