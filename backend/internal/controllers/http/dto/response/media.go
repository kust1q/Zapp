package response

type TweetMedia struct {
	MediaURL  string `json:"media_url"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}

type Avatar struct {
	AvatarURL string `json:"avatar_url"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}
