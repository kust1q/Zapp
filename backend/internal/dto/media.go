package dto

type TweetMediaResponse struct {
	ID        int    `json:"id"`
	TweetID   int    `json:"tweet_id"`
	MediaURL  string `json:"media_url"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}

type AvatarResponse struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	MediaURL  string `json:"media_url"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}
