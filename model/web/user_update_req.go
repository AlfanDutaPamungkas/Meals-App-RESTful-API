package web

type UserUpdateReq struct {
	Username string `json:"username" validate:"max=100"`
	Email    string `json:"email" validate:"omitempty,max=100,email"`
}
