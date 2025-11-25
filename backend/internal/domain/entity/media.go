package entity

type (
	TweetMedia struct {
		ID        int
		TweetID   int
		Path      string
		MimeType  string
		SizeBytes int64
	}

	Avatar struct {
		ID        int
		UserID    int
		Path      string
		MimeType  string
		SizeBytes int64
	}
)
