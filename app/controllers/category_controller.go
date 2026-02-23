package controllers

import (
	"errors"
	"net/http"
	"time"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/requests"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type CategoryController struct {
	service services.CategoryService
}

func NewCategoryController(service services.CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

// @Summary		List all categories
// @Description	Retrieve a list of all categories, including related products
// @Tags			categories
// @Security		Bearer
// @Produce		json
// @Success		200	{object}	helpers.ResponseParams[models.Category]{data=[]models.Category}
// @Failure		500	{object}	helpers.ResponseParams[any]
// @Router			/categories [get]
func (c *CategoryController) List(ctx *gin.Context) {
	categories, _, err := c.service.List(1, 1000)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to get categories list")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to get categories list",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Int("count", len(categories)).
		Msg("Categories list retrieved successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Category]{Data: &categories}, http.StatusOK)
}

// @Summary		Get a category by ID
// @Description	Retrieve a category by its ID, including related products
// @Tags			categories
// @Security		Bearer
// @Produce		json
// @Param			id	path		string	true	"Category ID"
// @Success		200	{object}	helpers.ResponseParams[models.Category]{item=models.Category}
// @Failure		404	{object}	helpers.ResponseParams[any]
// @Router			/categories/{id} [get]
func (c *CategoryController) Get(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	category, err := c.service.FindByID(id)
	if err != nil {
		log.Error().
			Err(err).
			Uint("category_id", id).
			Msg("Category not found")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": "Category not found"},
			Message:   "Category not found",
			Reference: "ERROR-2",
		}, http.StatusNotFound)
		return
	}

	log.Info().
		Uint("category_id", id).
		Str("category", category.Category).
		Msg("Category retrieved successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Category]{Item: category}, http.StatusOK)
}

// @Summary		Create a new category
// @Description	Create a new category
// @Tags			categories
// @Security		Bearer
// @Accept			json
// @Produce		json
// @Param			category	body		requests.CategoryRequest	true	"Category Data"
// @Success		201			{object}	helpers.ResponseParams[models.Category]{item=models.Category}
// @Failure		400			{object}	helpers.ResponseParams[any]
// @Failure		500			{object}	helpers.ResponseParams[any]
// @Router			/categories [post]
func (c *CategoryController) Create(ctx *gin.Context) {
	var req requests.CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			log.Warn().
				Err(err).
				Str("category", req.Category).
				Msg("Category creation validation failed")
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Errors:    helpers.ValidationError(verr),
				Message:   "Invalid parameters",
				Reference: "ERROR-4",
			}, http.StatusBadRequest)
			return
		}

		log.Warn().
			Err(err).
			Str("category", req.Category).
			Msg("Failed to bind category creation request")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to create category",
			Reference: "ERROR-3",
		}, http.StatusBadRequest)
		return
	}

	// Convert request to model
	category := models.Category{
		Category:  req.Category,
		UpdatedAt: time.Now(),
	}

	createdCategory, err := c.service.Create(category)
	if err != nil {
		log.Error().
			Err(err).
			Str("category", req.Category).
			Msg("Failed to create category")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to create category",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("category_id", createdCategory.ID).
		Str("category", createdCategory.Category).
		Msg("Category created successfully via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Category]{Item: createdCategory}, http.StatusCreated)
}

// @Summary		Update an existing category
// @Description	Update a category by its ID
// @Tags			categories
// @Security		Bearer
// @Accept			json
// @Produce		json
// @Param			id			path		uint						true	"Category ID"
// @Param			category	body		requests.CategoryRequest	true	"Category Data"
// @Success		200			{object}	helpers.ResponseParams[models.Category]{item=models.Category}
// @Failure		400			{object}	helpers.ResponseParams[any]
// @Failure		500			{object}	helpers.ResponseParams[any]
// @Router			/categories/{id} [put]
func (c *CategoryController) Update(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}

	var req requests.CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			log.Warn().
				Err(err).
				Uint("category_id", id).
				Msg("Category update validation failed")
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Errors:    helpers.ValidationError(verr),
				Message:   "Invalid parameters",
				Reference: "ERROR-4",
			}, http.StatusBadRequest)
			return
		}

		log.Warn().
			Err(err).
			Uint("category_id", id).
			Msg("Failed to bind category update request")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to update category",
			Reference: "ERROR-3",
		}, http.StatusBadRequest)
		return
	}

	// Convert request to model
	category := models.Category{
		Category:  req.Category,
		UpdatedAt: time.Now(),
	}

	updatedCategory, err := c.service.Update(id, category)
	if err != nil {
		log.Error().
			Err(err).
			Uint("category_id", id).
			Msg("Failed to update category")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to update category",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("category_id", updatedCategory.ID).
		Str("category", updatedCategory.Category).
		Msg("Category updated successfully via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Category]{Item: updatedCategory}, http.StatusOK)
}

// @Summary		Delete a category by ID
// @Description	Delete a specific category by its ID
// @Tags			categories
// @Security		Bearer
// @Produce		json
// @Param			id	path		string	true	"Category ID"
// @Success		200	{object}	helpers.ResponseParams[any]
// @Failure		404	{object}	helpers.ResponseParams[any]
// @Failure		500	{object}	helpers.ResponseParams[any]
// @Router			/categories/{id} [delete]
func (c *CategoryController) Delete(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	if err := c.service.DeleteByID(id); err != nil {
		log.Error().
			Err(err).
			Uint("category_id", id).
			Msg("Failed to delete category")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": "Category not found"},
			Message:   "Failed to delete category",
			Reference: "ERROR-2",
		}, http.StatusNotFound)
		return
	}

	log.Info().
		Uint("category_id", id).
		Msg("Category deleted successfully via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{Message: "Category deleted successfully"}, http.StatusOK)
}
