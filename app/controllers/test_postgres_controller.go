package controllers

import (
	"strconv"

	"golang_starter_kit_2025/app/helpers"
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
// @Tags         Test Postgres
// @Produce      json
// @Success      200  {array}   models.Test
// @Router       /tests [get]
func (c *TestController) List(ctx *gin.Context) {
	tests, err := c.service.GetAll()
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Gagal mendapatkan data test",
			Reference: "ERROR-TEST-1",
		}, 500)
		return
	}
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Test]{Data: &tests}, 200)
}

// Get godoc
// @Summary      Get test by ID
// @Description  Get a test record by ID from PostgreSQL
// @Tags         Test Postgres
// @Produce      json
// @Param        id   path      int  true  "Test ID"
// @Success      200  {object}  models.Test
// @Failure      404  {object}  map[string]string
// @Router       /tests/{id} [get]
func (c *TestController) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": "Invalid ID"},
			Message:   "ID tidak valid",
			Reference: "ERROR-TEST-2",
		}, 400)
		return
	}
	test, err := c.service.GetByID(uint(id))
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Data test tidak ditemukan",
			Reference: "ERROR-TEST-3",
		}, 404)
		return
	}
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Test]{Item: test}, 200)
}

// Create godoc
// @Summary      Create new test
// @Description  Create a new test record in PostgreSQL
// @Tags         Test Postgres
// @Accept       json
// @Produce      json
// @Param        test  body      models.Test  true  "Test Data"
// @Success      201   {object}  models.Test
// @Failure      400   {object}  map[string]string
// @Router       /tests [post]
func (c *TestController) Create(ctx *gin.Context) {
	var test models.Test
	if err := ctx.ShouldBindJSON(&test); err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Parameter tidak valid",
			Reference: "ERROR-TEST-4",
		}, 400)
		return
	}
	if err := c.service.Create(&test); err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Gagal membuat data test",
			Reference: "ERROR-TEST-5",
		}, 500)
		return
	}
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Test]{Item: &test}, 201)
}

// Update godoc
// @Summary      Update test
// @Description  Update a test record in PostgreSQL
// @Tags         Test Postgres
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
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": "Invalid ID"},
			Message:   "ID tidak valid",
			Reference: "ERROR-TEST-6",
		}, 400)
		return
	}
	var test models.Test
	if err := ctx.ShouldBindJSON(&test); err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Parameter tidak valid",
			Reference: "ERROR-TEST-7",
		}, 400)
		return
	}
	test.ID = uint(id)
	if err := c.service.Update(&test); err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Gagal update data test",
			Reference: "ERROR-TEST-8",
		}, 500)
		return
	}
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Test]{Item: &test}, 200)
}

// Delete godoc
// @Summary      Delete test
// @Description  Delete a test record in PostgreSQL
// @Tags         Test Postgres
// @Produce      json
// @Param        id   path      int  true  "Test ID"
// @Success      200  {object}  map[string]string
// @Router       /tests/{id} [delete]
func (c *TestController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": "Invalid ID"},
			Message:   "ID tidak valid",
			Reference: "ERROR-TEST-9",
		}, 400)
		return
	}
	if err := c.service.Delete(uint(id)); err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Gagal menghapus data test",
			Reference: "ERROR-TEST-10",
		}, 500)
		return
	}
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{}, 200)
}
