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

type RoleController struct {
	service services.RoleService
}

func NewRoleController(service services.RoleService) *RoleController {
	return &RoleController{service: service}
}

// @Summary		Get All Roles
// @Description	API to get all roles
// @Tags			Role
// @Accept			json
// @Produce		json
// @Success		200	{object}	helpers.ResponseParams[models.Role]{data=[]models.Role}
// @Router			/roles [get]
func (c *RoleController) List(ctx *gin.Context) {
	roles, _, err := c.service.List(1, 1000)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to get roles list")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to get roles list",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Int("count", len(roles)).
		Msg("Roles list retrieved successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Role]{Data: &roles}, http.StatusOK)
}

// @Summary		Create Role
// @Description	API to create a new role
// @Tags			Role
// @Accept			json
// @Produce		json
// @Param			role	body		requests.RoleRequestPut	true	"Role Data"
// @Success		201		{object}	helpers.ResponseParams[models.Role]{item=models.Role}
// @Router			/roles [post]
func (c *RoleController) Create(ctx *gin.Context) {
	var role models.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			log.Warn().
				Err(err).
				Str("name", role.Name).
				Msg("Role creation validation failed")
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Errors:    helpers.ValidationError(verr),
				Message:   "Invalid parameters",
				Reference: "ERROR-4",
			}, http.StatusBadRequest)
			return
		}

		log.Warn().
			Err(err).
			Str("name", role.Name).
			Msg("Failed to bind role creation request")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to create role",
			Reference: "ERROR-3",
		}, http.StatusBadRequest)
		return
	}

	createdRole, err := c.service.Create(role)
	if err != nil {
		log.Error().
			Err(err).
			Str("name", role.Name).
			Msg("Failed to create role")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to create role",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("role_id", createdRole.ID).
		Str("name", createdRole.Name).
		Msg("Role created successfully via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Role]{Item: createdRole}, http.StatusCreated)
}

// @Summary		Update Role
// @Description	API to update an existing role
// @Tags			Role
// @Accept			json
// @Produce		json
// @Param			id		path		uint					true	"Role ID"
// @Param			role	body		requests.RoleRequestPut	true	"Role Data"
// @Success		200		{object}	helpers.ResponseParams[models.Role]{item=models.Role}
// @Router			/roles/{id} [put]
func (c *RoleController) Update(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}

	var role models.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			log.Warn().
				Err(err).
				Uint("role_id", id).
				Msg("Role update validation failed")
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Errors:    helpers.ValidationError(verr),
				Message:   "Invalid parameters",
				Reference: "ERROR-4",
			}, http.StatusBadRequest)
			return
		}

		log.Warn().
			Err(err).
			Uint("role_id", id).
			Msg("Failed to bind role update request")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to update role",
			Reference: "ERROR-3",
		}, http.StatusBadRequest)
		return
	}

	updatedRole, err := c.service.Update(id, role)
	if err != nil {
		log.Error().
			Err(err).
			Uint("role_id", id).
			Msg("Failed to update role")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to update role",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("role_id", updatedRole.ID).
		Str("name", updatedRole.Name).
		Msg("Role updated successfully via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Role]{Item: updatedRole}, http.StatusOK)
}

// @Summary		Delete Role
// @Description	API to delete Role by ID
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
		log.Error().
			Err(err).
			Uint("role_id", id).
			Msg("Failed to delete role")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to delete role",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("role_id", id).
		Msg("Role deleted successfully via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{Message: "Role deleted successfully"}, http.StatusOK)
}

// Struct to wrap the permissions array
type AssignPermissionsRequest struct {
	Permissions []uint `json:"permissions"`
}

func (c *RoleController) AssignPermissions(ctx *gin.Context) {
	var req AssignPermissionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn().
			Err(err).
			Msg("Failed to bind assign permissions request")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Invalid parameters",
			Reference: "ERROR-4",
		}, http.StatusBadRequest)
		return
	}

	roleID, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	err := c.service.AssignPermissions(roleID, req.Permissions)
	if err != nil {
		log.Error().
			Err(err).
			Uint("role_id", roleID).
			Uints("permission_ids", req.Permissions).
			Msg("Failed to assign permissions to role")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to assign permissions to role",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("role_id", roleID).
		Uints("permission_ids", req.Permissions).
		Msg("Permissions assigned to role successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{Message: "Permissions assigned to role"}, http.StatusOK)
}

func (c *RoleController) GetPermissions(ctx *gin.Context) {
	roleID, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	permissions, err := c.service.GetPermissions(roleID)
	if err != nil {
		log.Error().
			Err(err).
			Uint("role_id", roleID).
			Msg("Failed to get role permissions")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to get role permissions",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("role_id", roleID).
		Int("permission_count", len(permissions)).
		Msg("Role permissions retrieved successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Permission]{Data: &permissions}, http.StatusOK)
}
