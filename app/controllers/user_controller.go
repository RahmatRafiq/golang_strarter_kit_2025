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
	users, err := c.service.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// @Summary	Show a user
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		id	path		string	true	"User ID"
// @Success	200	{object}	models.User
// @Router		/users/{id} [get]
func (c *UserController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.service.Find(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// @Summary	Upsert a user
// @Tags		users
// @Accept		json
// @Produce	json
// @Success	200	{object}	models.User
// @Router		/users [put]
// @Param		JSON	body	models.User	true	"User object"
func (c *UserController) Put(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedUser, err := c.service.Put(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, updatedUser)
}

// @Summary	Delete a user
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		id	path		string	true	"User ID"
// @Success	204	{string}	string	"User deleted"
// @Router		/users/{id} [delete]
func (c *UserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
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

	userId := ctx.Param("id")
	err := c.service.AssignRolesToUser(userId, req.Roles)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Roles assigned to user"})
}
func (c *UserController) GetRoles(ctx *gin.Context) {
	userId := ctx.Param("id")
	roles, err := c.service.GetRolesByUserId(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, roles)
}
