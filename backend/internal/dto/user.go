package dto

import "time"

type SmallUserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Bio       string    `json:"bio"`
	Gen       string    `json:"gen"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	AvatarURL string    `json:"avatar_url"`
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
