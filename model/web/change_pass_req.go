package web

type ChangePassReq struct {
	Password string `json:"password" validate:"required"`
}
