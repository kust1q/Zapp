package dto

import "time"

type CreateTweetRequest struct {
	Content string `json:"content" binding:"required,min=1,max=280"`
}

type TweetResponse struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ParentTweetID *int      `json:"parent_tweet_id,omitempty"`
	ReplyCount    int       `json:"reply_count"`
	RetweetCount  int       `json:"retweet_count"`
	LikeCount     int       `json:"like_count"`
}

type UpdateTweetRequest struct {
	Content string `json:"content" binding:"required,min=1,max=280"`
}

type UpdateTweetResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}
