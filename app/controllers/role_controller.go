package controllers

import (
	"errors"
	"net/http"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RoleController struct {
	service services.RoleService
}

func NewRoleController(service services.RoleService) *RoleController {
	return &RoleController{service: service}
}

// @Summary		Get All Roles
// @Description	API untuk mendapatkan semua Role
// @Tags			Role
// @Accept			json
// @Produce		json
// @Success		200	{object}	helpers.ResponseParams[models.Role]{data=[]models.Role}
// @Router			/roles [get]
func (c *RoleController) List(ctx *gin.Context) {
	roles, err := c.service.GetAll()
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to get roles list",
			Reference: "ERROR-3",
		}, 500)
		return
	}

	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Role]{Data: &roles}, 200)
}

// @Summary		Create/Update Role
// @Description	API untuk mengupdate atau membuat Role
// @Tags			Role
// @Accept			json
// @Produce		json
// @Param			role	body		requests.RoleRequestPut	true	"Role Data"
// @Success		200		{object}	helpers.ResponseParams[models.Role]{item=models.Role}
// @Router			/roles [put]
func (c *RoleController) Put(ctx *gin.Context) {
	var role models.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Errors:    helpers.ValidationError(verr),
				Message:   "Invalid parameters",
				Reference: "ERROR-4",
			}, 400)
		}

		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to create role",
			Reference: "ERROR-3",
		}, 400)
		return
	}
	updatedRole, err := c.service.Put(role)
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to create role",
			Reference: "ERROR-3",
		}, 400)
		return
	}

	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Role]{Item: &updatedRole}, 200)
}

// @Summary		Delete Role
// @Description	API untuk menghapus Role
// @Tags			Role
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"Role ID"
// @Success		200	{object}	helpers.ResponseParams[any]{}
// @Router			/roles/{id} [delete]
func (c *RoleController) Delete(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	if err := c.service.DeleteByID(id); err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to delete role",
			Reference: "ERROR-3",
		}, 500)
		return
	}

	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{}, 200)
}

// Struct to wrap the permissions array
type AssignPermissionsRequest struct {
	Permissions []uint `json:"permissions"`
}

func (c *RoleController) AssignPermissions(ctx *gin.Context) {
	var req AssignPermissionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roleID, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	err := c.service.AssignPermissions(roleID, req.Permissions)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Permissions assigned to role"})
}

func (c *RoleController) GetPermissions(ctx *gin.Context) {
	roleID, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	permissions, err := c.service.GetPermissions(roleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, permissions)
}
