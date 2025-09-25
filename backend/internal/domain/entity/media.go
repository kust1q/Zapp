package entity

type TweetMedia struct {
	ID        int    `db:"id"`
	TweetID   int    `db:"tweet_id"`
	Path      string `db:"path"`
	MimeType  string `db:"mime_type"`
	SizeBytes int64  `db:"size_bytes"`
}

type Avatar struct {
	ID        int    `db:"id"`
	UserID    int    `db:"user_id"`
	Path      string `db:"path"`
	MimeType  string `db:"mime_type"`
	SizeBytes int64  `db:"size_bytes"`
}
