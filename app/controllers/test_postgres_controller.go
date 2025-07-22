package controllers

import (
	"net/http"
	"strconv"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
)

type TestController struct {
	service services.TestService
}

func NewTestController(service services.TestService) *TestController {
	return &TestController{service: service}
}

// List godoc
// @Summary      Get all test data
// @Description  Get all test records from PostgreSQL
// @Tags         Test
// @Produce      json
// @Success      200  {array}   models.Test
// @Router       /tests [get]
func (c *TestController) List(ctx *gin.Context) {
	tests, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tests)
}

// Get godoc
// @Summary      Get test by ID
// @Description  Get a test record by ID from PostgreSQL
// @Tags         Test
// @Produce      json
// @Param        id   path      int  true  "Test ID"
// @Success      200  {object}  models.Test
// @Failure      404  {object}  map[string]string
// @Router       /tests/{id} [get]
func (c *TestController) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	test, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, test)
}

// Create godoc
// @Summary      Create new test
// @Description  Create a new test record in PostgreSQL
// @Tags         Test
// @Accept       json
// @Produce      json
// @Param        test  body      models.Test  true  "Test Data"
// @Success      201   {object}  models.Test
// @Failure      400   {object}  map[string]string
// @Router       /tests [post]
func (c *TestController) Create(ctx *gin.Context) {
	var test models.Test
	if err := ctx.ShouldBindJSON(&test); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.service.Create(&test); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, test)
}

// Update godoc
// @Summary      Update test
// @Description  Update a test record in PostgreSQL
// @Tags         Test
// @Accept       json
// @Produce      json
// @Param        id    path      int         true  "Test ID"
// @Param        test  body      models.Test true  "Test Data"
// @Success      200   {object}  models.Test
// @Failure      400   {object}  map[string]string
// @Router       /tests/{id} [put]
func (c *TestController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var test models.Test
	if err := ctx.ShouldBindJSON(&test); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	test.ID = uint(id)
	if err := c.service.Update(&test); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, test)
}

// Delete godoc
// @Summary      Delete test
// @Description  Delete a test record in PostgreSQL
// @Tags         Test
// @Produce      json
// @Param        id   path      int  true  "Test ID"
// @Success      200  {object}  map[string]string
// @Router       /tests/{id} [delete]
func (c *TestController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
