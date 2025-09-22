package dto

type TweetMedia struct {
	ID        int    `json:"id" db:"id"`
	TweetID   int    `json:"tweet_id" db:"tweet_id"`
	MediaURL  string `json:"media_url" db:"media_url"`
	MimeType  string `json:"mime_type" db:"mime_type"`
	SizeBytes int64  `json:"size_bytes" db:"size_bytes"`
}

type Avatar struct {
	ID        int    `json:"id" db:"id"`
	UserID    int    `json:"user_id" db:"user_id"`
	MediaURL  string `json:"media_url" db:"media_url"`
	MimeType  string `json:"mime_type" db:"mime_type"`
	SizeBytes int64  `json:"size_bytes" db:"size_bytes"`
}
