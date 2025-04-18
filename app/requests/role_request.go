package requests

type RoleRequestPut struct {
	ID    uint   `json:"id" form:"id"`
	Name  string `json:"name" form:"name" binding:"required" example:"Admin" validate:"required"`
	Group string `json:"group" form:"group" binding:"required" example:"User" validate:"required"`
}

type RoleRequestAssignPermissions struct {
	PermissionIDs []uint `json:"permissions" form:"permissions" binding:"required" validate:"required"`
}
