package requests

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required,min=64,max=64"`
}
