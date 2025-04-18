package requests

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@mail.com"`
	Password string `json:"password" binding:"required" example:"12345678"`
}