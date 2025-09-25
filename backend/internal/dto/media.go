package dto

type TweetMediaDataResponse struct {
	ID        int    `json:"id"`
	TweetID   int    `json:"tweet_id"`
	MediaURL  string `json:"avatar_url"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}

type AvatarDataResponse struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	AvatarURL string `json:"avatar_url"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}
