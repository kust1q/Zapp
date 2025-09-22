package dto

import "time"

type SignUpRequest struct {
	Username       string `json:"username" binding:"required,alphanum,max=64"`
	Email          string `json:"email" binding:"required,email,max=64"`
	Password       string `json:"password" binding:"required,alphanum,min=8,max=64"`
	Gen            string `json:"gen" binding:"required,oneof=male female"`
	Bio            string `json:"bio"`
	SecretQuestion string `json:"secret_question" binding:"required,min=8,max=100"`
	SecretAnswer   string `json:"secret_answer" binding:"required,min=8,max=100"`
}

type SignUpResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Bio       string    `json:"bio"`
	Gen       string    `json:"gen"`
	Avatar    Avatar    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,alphanum,min=8,max=64"`
}

type RefreshRequest struct {
	Refresh string `json:"refresh" binding:"required"`
}

type SignInResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type UpdateSecuritySettingsRequest struct {
	OldSecretAnswer   string `json:"old_secret_answer" binding:"required"`
	NewSecretQuestion string `json:"new_secret_question" binding:"required,min=8,max=100"`
	NewSecretAnswer   string `json:"new_secret_answer" binding:"required,min=8,max=50"`
}

type ResetPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,alphanum,min=8,max=64"`
	NewPassword string `json:"new_password" binding:"required,alphanum,min=8,max=64"`
}

type ForgotPasswordRequest struct {
	Email        string `json:"email" binding:"required,email,max=64"`
	SecretAnswer string `json:"secret_answer" binding:"required,min=8,max=50"`
	NewPassword  string `json:"new_password" binding:"required,alphanum,min=8,max=64"`
}
