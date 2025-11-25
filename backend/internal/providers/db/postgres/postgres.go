package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kust1q/Zapp/backend/internal/config"
)

const (
	UserTable           = "users"
	FollowsTable        = "follows"
	LikesTable          = "likes"
	TweetsTable         = "tweets"
	RetweetsTable       = "retweets"
	TweetMediaTable     = "tweet_media"
	SecretQuestionTable = "secret_questions"
	AvatarsTable        = "avatars"
)

type PostgresDB struct {
	db        *sqlx.DB
	userCache UserCache
}

func NewPostgresDB(cfg config.PostgresConfig, usercCache UserCache) (*PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresDB{
		db:        db,
		userCache: usercCache,
	}, nil
}

func (pg *PostgresDB) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return pg.db.BeginTx(ctx, nil)
}
