package dto

import "time"

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum,max=64"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,alphanum,min=8,max=64"`
	Gen      string `json:"gen" binding:"required,oneof=male female"`
	Bio      string `json:"bio"`
}

type UpdateUserRequest struct {
	Password string `json:"password" binding:"required,alphanum,min=8,max=64"`
	Bio      string `json:"bio"`
}

type UserResponse struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Bio         string    `json:"bio"`
	Gen         string    `json:"gen"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	IsSuperuser bool      `json:"is_superuser"`
}
