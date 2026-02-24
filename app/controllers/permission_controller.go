package controllers

import (
	"errors"
	"net/http"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
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
		log.Error().
			Err(err).
			Msg("Failed to get permissions list")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to get permissions list",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Int("count", len(permissions)).
		Msg("Permissions list retrieved successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Permission]{Data: &permissions}, http.StatusOK)
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
			log.Warn().
				Err(err).
				Str("name", permission.Name).
				Msg("Permission creation validation failed")
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Errors:    helpers.ValidationError(verr),
				Message:   "Invalid parameters",
				Reference: "ERROR-4",
			}, http.StatusBadRequest)
			return
		}
		log.Warn().
			Err(err).
			Str("name", permission.Name).
			Msg("Failed to bind permission creation request")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to create permission",
			Reference: "ERROR-3",
		}, http.StatusBadRequest)
		return
	}

	createdPermission, err := c.service.Create(permission)
	if err != nil {
		log.Error().
			Err(err).
			Str("name", permission.Name).
			Msg("Failed to create permission")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to create permission",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("permission_id", createdPermission.ID).
		Str("name", createdPermission.Name).
		Msg("Permission created successfully via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Permission]{Item: createdPermission}, http.StatusCreated)
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
			log.Warn().
				Err(err).
				Uint("permission_id", id).
				Msg("Permission update validation failed")
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Errors:    helpers.ValidationError(verr),
				Message:   "Invalid parameters",
				Reference: "ERROR-4",
			}, http.StatusBadRequest)
			return
		}
		log.Warn().
			Err(err).
			Uint("permission_id", id).
			Msg("Failed to bind permission update request")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to update permission",
			Reference: "ERROR-3",
		}, http.StatusBadRequest)
		return
	}

	updatedPermission, err := c.service.Update(id, permission)
	if err != nil {
		log.Error().
			Err(err).
			Uint("permission_id", id).
			Msg("Failed to update permission")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to update permission",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("permission_id", updatedPermission.ID).
		Str("name", updatedPermission.Name).
		Msg("Permission updated successfully via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Permission]{Item: updatedPermission}, http.StatusOK)
}

// @Summary		Delete Permission
// @Description	API to delete Permission by ID
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
		log.Error().
			Err(err).
			Uint("permission_id", id).
			Msg("Failed to delete permission")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to delete permission",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("permission_id", id).
		Msg("Permission deleted successfully via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Permission]{Message: "Permission deleted successfully"}, http.StatusOK)
}
