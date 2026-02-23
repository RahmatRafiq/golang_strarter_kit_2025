package controllers

import (
	"net/http"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
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
// @Success	200	{array}	models.User
// @Router		/users [get]
func (c *UserController) List(ctx *gin.Context) {
	// FIXED: Use List() with pagination instead of removed GetAllUsers()
	// Default to page 1, limit 100 for backward compatibility
	users, total, err := c.service.List(1, 100)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
		"page":  1,
		"limit": 100,
	})
}

// @Summary	Show a user
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		id	path		string	true	"User ID"
// @Success	200	{object}	models.User
// @Router		/users/{id} [get]
func (c *UserController) Get(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	user, err := c.service.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// @Summary	Create a new user
// @Tags		users
// @Accept		json
// @Produce	json
// @Success	201	{object}	models.User
// @Router		/users [post]
// @Param		JSON	body	models.User	true	"User object"
func (c *UserController) Create(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdUser, err := c.service.Create(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdUser)
}

// @Summary	Update an existing user
// @Tags		users
// @Accept		json
// @Produce	json
// @Success	200	{object}	models.User
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := c.service.Update(id, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedUser)
}

// @Summary	Delete a user
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		id	path		string	true	"User ID"
// @Success	204	{string}	string	"User deleted"
// @Router		/users/{id} [delete]
func (c *UserController) Delete(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	if err := c.service.DeleteByID(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// Struct to wrap the roles array
type AssignRolesRequest struct {
	Roles []uint `json:"roles"`
}

func (c *UserController) AssignRoles(ctx *gin.Context) {
	var req AssignRolesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	err := c.service.AssignRoles(userID, req.Roles)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Roles assigned to user"})
}
func (c *UserController) GetRoles(ctx *gin.Context) {
	userID, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	roles, err := c.service.GetRoles(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, roles)
}
