package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kust1q/Zapp/backend/internal/config"
)

const (
	UserTable       = "users"
	FollowsTable    = "follows"
	LikesTable      = "likes"
	TweetsTable     = "tweets"
	RetweetsTable   = "retweets"
	TweetMediaTable = "tweet_media"
)

func NewPostgresDB(cfg config.PostgresConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
