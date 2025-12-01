package entity

type (
	MediaType string

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

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	MediaTypeAudio MediaType = "audio"
	MediaTypeGIF   MediaType = "gif"
)
