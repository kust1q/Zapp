package entity

type TweetMedia struct {
	ID        int    `db:"id"`
	TweetID   int    `db:"tweet_id"`
	MediaURL  string `db:"media_url"`
	MimeType  string `db:"mime_type"`
	SizeBytes int64  `db:"size_bytes"`
}

type Avatar struct {
	ID        int    `db:"id"`
	UserID    int    `db:"user_id"`
	MediaURL  string `db:"media_url"`
	MimeType  string `db:"mime_type"`
	SizeBytes int64  `db:"size_bytes"`
}
