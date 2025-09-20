package entity

import (
	"time"
)

type Tweet struct {
	ID            int       `db:"id"`
	UserID        int       `db:"user_id"`
	ParentTweetID int       `db:"parent_tweet_id"`
	Content       string    `db:"content"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type Retweet struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	TweetID   int       `db:"tweet_id"`
	CreatedAt time.Time `db:"created_at"`
}
