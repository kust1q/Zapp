package dto

import "time"

type SmallUserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Avatar   Avatar `json:"avatar"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Bio       string    `json:"bio"`
	Gen       string    `json:"gen"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Avatar    Avatar    `json:"avatar"`
}

type UserProfileResponse struct {
	UserResponse
	Tweets []TweetResponse `json:"tweets"`
}

type UpdateBioRequest struct {
	Bio string `json:"bio"`
}

type FollowResponse struct {
	FollowerID  int       `json:"follower_id" db:"follower_id"`
	FollowingID int       `json:"following_id" db:"following_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type AvatarData struct {
	MediaURL  string `json:"media_url"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}
