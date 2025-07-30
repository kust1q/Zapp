package dto

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Gen      string `json:"gen" binding:"required,oneof=Male Female"`
	Bio      string `json:"bio"`
}

type UpdateUserRequest struct {
}
