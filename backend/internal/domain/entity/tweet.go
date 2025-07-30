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
	ReplyCount    int       `db:"reply_count"`
	RetweetCount  int       `db:"retweet_count"`
	LikeCount     int       `db:"like_count"`
}

type TweetMedia struct {
	TweetID   int    `db:"tweet_id"`
	MediaURL  string `db:"media_url"`
	MediaType string `db:"media_type"`
}

type Like struct {
	UserID  int `db:"user_id"`
	TweetID int `db:"tweet_id"`
}

type Retweet struct {
	UserID    int       `db:"user_id"`
	TweetID   int       `db:"tweet_id"`
	CreatedAt time.Time `db:"created_at"`
}
