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

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

// @Scheme
// @Summary	Show all users
// @Tags		users
// @Accept		json
// @Produce	json
// @Success	200	{object}	helpers.ResponseParams[models.User]{data=[]models.User}
// @Router		/users [get]
func (c *UserController) List(ctx *gin.Context) {
	// FIXED: Use List() with pagination instead of removed GetAllUsers()
	// Default to page 1, limit 100 for backward compatibility
	users, total, err := c.service.List(1, 100)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to get users list")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to get users list",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Int("count", len(users)).
		Int64("total", total).
		Msg("Users list retrieved successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.User]{
		Data:    &users,
		Message: "Users list retrieved successfully",
	}, http.StatusOK)
}

// @Summary	Show a user
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		id	path		string	true	"User ID"
// @Success	200	{object}	helpers.ResponseParams[models.User]{item=models.User}
// @Router		/users/{id} [get]
func (c *UserController) Get(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	user, err := c.service.FindByID(id)
	if err != nil {
		log.Error().
			Err(err).
			Uint("user_id", id).
			Msg("User not found")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": "User not found"},
			Message:   "User not found",
			Reference: "ERROR-2",
		}, http.StatusNotFound)
		return
	}

	log.Info().
		Uint("user_id", id).
		Str("email", user.Email).
		Msg("User retrieved successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.User]{Item: user}, http.StatusOK)
}

// @Summary	Create a new user
// @Tags		users
// @Accept		json
// @Produce	json
// @Success	201	{object}	helpers.ResponseParams[models.User]{item=models.User}
// @Router		/users [post]
// @Param		JSON	body	models.User	true	"User object"
func (c *UserController) Create(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			log.Warn().
				Err(err).
				Str("email", user.Email).
				Msg("User creation validation failed")
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Errors:    helpers.ValidationError(verr),
				Message:   "Invalid parameters",
				Reference: "ERROR-4",
			}, http.StatusBadRequest)
			return
		}

		log.Warn().
			Err(err).
			Str("email", user.Email).
			Msg("Failed to bind user creation request")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to create user",
			Reference: "ERROR-3",
		}, http.StatusBadRequest)
		return
	}

	createdUser, err := c.service.Create(user)
	if err != nil {
		log.Error().
			Err(err).
			Str("email", user.Email).
			Msg("Failed to create user")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to create user",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("user_id", createdUser.ID).
		Str("email", createdUser.Email).
		Msg("User created successfully via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.User]{Item: createdUser}, http.StatusCreated)
}

// @Summary	Update an existing user
// @Tags		users
// @Accept		json
// @Produce	json
// @Success	200	{object}	helpers.ResponseParams[models.User]{item=models.User}
// @Router		/users/{id} [put]
// @Param		id		path	uint		true	"User ID"
// @Param		JSON	body	models.User	true	"User object"
func (c *UserController) Update(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}

	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			log.Warn().
				Err(err).
				Uint("user_id", id).
				Msg("User update validation failed")
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Errors:    helpers.ValidationError(verr),
				Message:   "Invalid parameters",
				Reference: "ERROR-4",
			}, http.StatusBadRequest)
			return
		}

		log.Warn().
			Err(err).
			Uint("user_id", id).
			Msg("Failed to bind user update request")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to update user",
			Reference: "ERROR-3",
		}, http.StatusBadRequest)
		return
	}

	updatedUser, err := c.service.Update(id, user)
	if err != nil {
		log.Error().
			Err(err).
			Uint("user_id", id).
			Msg("Failed to update user")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to update user",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("user_id", updatedUser.ID).
		Str("email", updatedUser.Email).
		Msg("User updated successfully via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.User]{Item: updatedUser}, http.StatusOK)
}

// @Summary	Delete a user
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		id	path		string	true	"User ID"
// @Success	200	{object}	helpers.ResponseParams[any]
// @Router		/users/{id} [delete]
func (c *UserController) Delete(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	if err := c.service.DeleteByID(id); err != nil {
		log.Error().
			Err(err).
			Uint("user_id", id).
			Msg("Failed to delete user")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": "User not found"},
			Message:   "Failed to delete user",
			Reference: "ERROR-2",
		}, http.StatusNotFound)
		return
	}

	log.Info().
		Uint("user_id", id).
		Msg("User deleted successfully via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{Message: "User deleted successfully"}, http.StatusOK)
}

// Struct to wrap the roles array
type AssignRolesRequest struct {
	Roles []uint `json:"roles"`
}

func (c *UserController) AssignRoles(ctx *gin.Context) {
	var req AssignRolesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn().
			Err(err).
			Msg("Failed to bind assign roles request")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Invalid parameters",
			Reference: "ERROR-4",
		}, http.StatusBadRequest)
		return
	}

	userID, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	err := c.service.AssignRoles(userID, req.Roles)
	if err != nil {
		log.Error().
			Err(err).
			Uint("user_id", userID).
			Uints("role_ids", req.Roles).
			Msg("Failed to assign roles to user")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to assign roles to user",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("user_id", userID).
		Uints("role_ids", req.Roles).
		Msg("Roles assigned to user successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{Message: "Roles assigned to user"}, http.StatusOK)
}
func (c *UserController) GetRoles(ctx *gin.Context) {
	userID, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	roles, err := c.service.GetRoles(userID)
	if err != nil {
		log.Error().
			Err(err).
			Uint("user_id", userID).
			Msg("Failed to get user roles")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to get user roles",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("user_id", userID).
		Int("role_count", len(roles)).
		Msg("User roles retrieved successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Role]{Data: &roles}, http.StatusOK)
}
