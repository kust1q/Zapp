package dto

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,alphanum,min=8,max=64"`
}

type SignInResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
