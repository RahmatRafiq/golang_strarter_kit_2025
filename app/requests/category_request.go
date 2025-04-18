package requests

type CategoryRequest struct {
	ID       uint   `json:"id,omitempty" form:"id,omitempty"`
	Category string `json:"category" form:"category" binding:"required" example:"Electronics" validate:"required"`
}
