package web

type UserRegisterReq struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=100,email"`
	Password string `json:"password" validate:"required"`
}
