package controllers

import (
	"errors"
	"net/http"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/requests"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type ProductController struct {
	service services.ProductService
}

func NewProductController(service services.ProductService) *ProductController {
	return &ProductController{
		service: service,
	}
}

// @Summary		Get all products
// @Description	API untuk mendapatkan semua produk
// @Tags			Product
// @Accept			json
// @Produce		json
// @Param			request	query		requests.FilterRequest	false	"Filter request"
// @Success		200		{object}	helpers.ResponseParams[models.Product]{data=[]models.Product}
// @Router			/products [get]
func (c *ProductController) List(ctx *gin.Context) {
	var filters requests.FilterRequest
	if err := ctx.ShouldBindQuery(&filters); err != nil {
		log.Warn().
			Err(err).
			Msg("Invalid filter parameters for products list")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "Invalid parameters",
			Reference: "ERROR-4",
		}, http.StatusBadRequest)
		return
	}

	products, _, err := c.service.List(1, 1000)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to get products list")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "Failed to get products list",
			Reference: "ERROR-3",
			Errors:    map[string]string{"error": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Int("count", len(products)).
		Msg("Products list retrieved successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Product]{Data: &products}, http.StatusOK)
}

// @Summary		Get product by ID
// @Description	API untuk mendapatkan produk berdasarkan ID
// @Tags			Product
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Product ID"
// @Success		200	{object}	helpers.ResponseParams[models.Product]{item=models.Product}
// @Router			/products/{id} [get]
func (c *ProductController) Get(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	product, err := c.service.FindByID(id)
	if err != nil {
		log.Error().
			Err(err).
			Uint("product_id", id).
			Msg("Product not found")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "Failed to get product",
			Reference: "ERROR-2",
			Errors:    map[string]string{"error": err.Error()},
		}, http.StatusNotFound)
		return
	}

	log.Info().
		Uint("product_id", id).
		Str("product_name", product.Name).
		Msg("Product retrieved successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Product]{Item: product}, http.StatusOK)
}

// @Summary		Create product
// @Description	API untuk membuat produk baru
// @Tags			Product
// @Accept			json
// @Produce		json
// @Param			product	body		requests.ProductRequest	true	"Product request body"
// @Success		201		{object}	helpers.ResponseParams[models.Product]{item=models.Product}
// @Router			/products [post]
func (c *ProductController) Create(ctx *gin.Context) {
	var request requests.ProductRequest
	if err := ctx.ShouldBind(&request); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			log.Warn().
				Err(err).
				Str("product_name", request.Name).
				Msg("Product creation validation failed")
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Errors:    helpers.ValidationError(verr),
				Message:   "Invalid parameters",
				Reference: "ERROR-4",
			}, http.StatusBadRequest)
			return
		}
	}

	product, err := c.service.Create(ctx, request)
	if err != nil {
		log.Error().
			Err(err).
			Str("product_name", request.Name).
			Msg("Failed to create product")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "Failed to create product",
			Reference: "ERROR-3",
			Errors:    map[string]string{"error": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("product_id", product.ID).
		Str("product_name", product.Name).
		Msg("Product created successfully via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Product]{Item: product}, http.StatusCreated)
}

// @Summary		Update product
// @Description	API untuk mengupdate produk yang sudah ada
// @Tags			Product
// @Accept			json
// @Produce		json
// @Param			id		path		uint					true	"Product ID"
// @Param			product	body		requests.ProductRequest	true	"Product request body"
// @Success		200		{object}	helpers.ResponseParams[models.Product]{item=models.Product}
// @Router			/products/{id} [put]
func (c *ProductController) Update(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}

	var request requests.ProductRequest
	if err := ctx.ShouldBind(&request); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			log.Warn().
				Err(err).
				Uint("product_id", id).
				Msg("Product update validation failed")
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Errors:    helpers.ValidationError(verr),
				Message:   "Invalid parameters",
				Reference: "ERROR-4",
			}, http.StatusBadRequest)
			return
		}
	}

	product, err := c.service.Update(ctx, id, request)
	if err != nil {
		log.Error().
			Err(err).
			Uint("product_id", id).
			Msg("Failed to update product")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "Failed to update product",
			Reference: "ERROR-3",
			Errors:    map[string]string{"error": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

	log.Info().
		Uint("product_id", product.ID).
		Str("product_name", product.Name).
		Msg("Product updated successfully via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Product]{Item: product}, http.StatusOK)
}

// @Summary		Delete product
// @Description	API to delete product by ID
// @Tags			Product
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Product ID"
// @Success		200	{object}	helpers.ResponseParams[models.Product]
// @Router			/products/{id} [delete]
func (c *ProductController) Delete(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	if err := c.service.DeleteByID(id); err != nil {
		log.Error().
			Err(err).
			Uint("product_id", id).
			Msg("Failed to delete product")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "Failed to delete product",
			Reference: "ERROR-3",
			Errors:    map[string]string{"error": err.Error()},
		}, http.StatusNotFound)
		return
	}

	log.Info().
		Uint("product_id", id).
		Msg("Product deleted successfully via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{Message: "Product deleted successfully"}, http.StatusOK)
}
