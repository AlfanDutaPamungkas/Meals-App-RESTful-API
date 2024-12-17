package web

type LoginRequest struct{
	Email    string `json:"email" validate:"required,max=100,email"`
	Password string `json:"password" validate:"required"`
}