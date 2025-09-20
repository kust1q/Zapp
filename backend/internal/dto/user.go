package dto

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Avatar   Avatar `json:"avatar"`
}

type Avatar struct {
	MediaURL  string `json:"media_url"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}
