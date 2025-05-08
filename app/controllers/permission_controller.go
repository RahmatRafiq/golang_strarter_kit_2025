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
	permissions, err := c.service.GetAll()
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Gagal mendapatkan daftar Permission",
			Reference: "ERROR-3",
		}, 500)
		return
	}

	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Permission]{Data: &permissions}, 200)
}

// @Summary		Create/Update Permission
// @Description	API untuk mengupdate atau membuat Permission
// @Tags			Permission
// @Accept			json
// @Produce		json
// @Param			permission	body		requests.PermissionRequest	true	"Permission Data"
// @Success		200			{object}	helpers.ResponseParams[models.Permission]{item=models.Permission}
// @Router			/permissions [put]
func (c *PermissionController) Put(ctx *gin.Context) {
	var permission models.Permission
	if err := ctx.ShouldBindJSON(&permission); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Errors:    helpers.ValidationError(verr),
				Message:   "Parameter tidak valid",
				Reference: "ERROR-4",
			}, 400)
			return
		}
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Gagal membuat Permission",
			Reference: "ERROR-3",
		}, 400)
		return
	}

	updatedPermission, err := c.service.Put(permission)
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Gagal membuat Permission",
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
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Gagal menghapus Permission",
			Reference: "ERROR-3",
		}, 400)
		return
	}

	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Permission]{Message: "Permission deleted"}, 200)
}
