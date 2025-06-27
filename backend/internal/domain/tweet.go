package domain

import "time"

type Tweet struct {
	Id              int       `json:"-"`
	User_id         int       `json:"user_id"`
	Parent_tweet_id int       `json:"parent_tweet_id"`
	Content         string    `json:"content"`
	Created_at      time.Time `json:"created_at"`
	Reply_count     int       `json:"reply_count"`
	Retweet_count   int       `json:"retweet_count"`
	Like_count      int       `json:"like_count"`
}

type Tweet_media struct {
	Tweet_id   int    `json:"tweet_id"`
	Media_url  string `json:"media_url"`
	Media_type string `json:"media_type"`
}

type Like struct {
	User_id  int `json:"user_id"`
	Tweet_id int `json:"tweet_id"`
}

type Retweet struct {
	User_id    int       `json:"user_id"`
	Tweet_id   int       `json:"tweet_id"`
	Created_at time.Time `json:"created_at"`
}
