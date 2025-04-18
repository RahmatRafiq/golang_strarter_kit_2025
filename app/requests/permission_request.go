package requests

type PermissionRequest struct {
	ID    uint   `json:"id" form:"id"`
	Name  string `json:"name" form:"name" binding:"required" example:"Create User" validate:"required"`
	Group string `json:"group" form:"group" binding:"required" example:"User" validate:"required"`
}
