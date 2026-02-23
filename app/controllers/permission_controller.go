package controllers

import (
	"errors"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PermissionController struct {
	service services.PermissionService
}

func NewPermissionController(service services.PermissionService) *PermissionController {
	return &PermissionController{service: service}
}

// @Summary		Get All Permissions
// @Description	API untuk mendapatkan semua Permission
// @Tags			Permission
// @Accept			json
// @Produce		json
// @Success		200	{object}	helpers.ResponseParams[models.Permission]{data=[]models.Permission}
// @Router			/permissions [get]
func (c *PermissionController) List(ctx *gin.Context) {
	permissions, _, err := c.service.List(1, 1000)
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to get permissions list",
			Reference: "ERROR-3",
		}, 500)
		return
	}

	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Permission]{Data: &permissions}, 200)
}

// @Summary		Create Permission
// @Description	API untuk membuat Permission baru
// @Tags			Permission
// @Accept			json
// @Produce		json
// @Param			permission	body		requests.PermissionRequest	true	"Permission Data"
// @Success		201			{object}	helpers.ResponseParams[models.Permission]{item=models.Permission}
// @Router			/permissions [post]
func (c *PermissionController) Create(ctx *gin.Context) {
	var permission models.Permission
	if err := ctx.ShouldBindJSON(&permission); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Errors:    helpers.ValidationError(verr),
				Message:   "Invalid parameters",
				Reference: "ERROR-4",
			}, 400)
			return
		}
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to create permission",
			Reference: "ERROR-3",
		}, 400)
		return
	}

	createdPermission, err := c.service.Create(permission)
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to create permission",
			Reference: "ERROR-3",
		}, 400)
		return
	}

	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Permission]{Item: &createdPermission}, 201)
}

// @Summary		Update Permission
// @Description	API untuk mengupdate Permission yang sudah ada
// @Tags			Permission
// @Accept			json
// @Produce		json
// @Param			id			path		uint						true	"Permission ID"
// @Param			permission	body		requests.PermissionRequest	true	"Permission Data"
// @Success		200			{object}	helpers.ResponseParams[models.Permission]{item=models.Permission}
// @Router			/permissions/{id} [put]
func (c *PermissionController) Update(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}

	var permission models.Permission
	if err := ctx.ShouldBindJSON(&permission); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Errors:    helpers.ValidationError(verr),
				Message:   "Invalid parameters",
				Reference: "ERROR-4",
			}, 400)
			return
		}
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to update permission",
			Reference: "ERROR-3",
		}, 400)
		return
	}

	updatedPermission, err := c.service.Update(id, permission)
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to update permission",
			Reference: "ERROR-3",
		}, 400)
		return
	}

	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Permission]{Item: &updatedPermission}, 200)
}

// @Summary		Delete Permission
// @Description	API untuk menghapus Permission berdasarkan ID
// @Tags			Permission
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"Permission ID"
// @Success		200	{object}	helpers.ResponseParams[models.Permission]{}
// @Router			/permissions/{id} [delete]
func (c *PermissionController) Delete(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	if err := c.service.DeleteByID(id); err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to delete permission",
			Reference: "ERROR-3",
		}, 400)
		return
	}

	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Permission]{Message: "Permission deleted"}, 200)
}
