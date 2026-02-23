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
func (c *ProductController) GetAll(ctx *gin.Context) {
	var filters requests.FilterRequest
	if err := ctx.ShouldBindQuery(&filters); err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "Invalid parameters",
			Reference: "ERROR-4",
		}, http.StatusBadRequest)
		return
	}

	products, err := c.service.GetAll(filters)
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "Failed to get products list",
			Reference: "ERROR-3",
			Errors:    map[string]string{"error": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

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
func (c *ProductController) GetByID(ctx *gin.Context) {
	id, ok := ParseIDParam(ctx)
	if !ok {
		return
	}
	product, err := c.service.FindByID(id)
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "Failed to get product",
			Reference: "ERROR-2",
			Errors:    map[string]string{"error": err.Error()},
		}, http.StatusNotFound)
		return
	}

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
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Message:   "Periksa kembali form anda",
				Errors:    helpers.ValidationError(verr),
				Reference: "ERROR-4",
			}, http.StatusBadRequest)
			return
		}
	}

	product, err := c.service.CreateProduct(ctx, request)
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "Failed to create product",
			Reference: "ERROR-3",
			Errors:    map[string]string{"error": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

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
			helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
				Message:   "Periksa kembali form anda",
				Errors:    helpers.ValidationError(verr),
				Reference: "ERROR-4",
			}, http.StatusBadRequest)
			return
		}
	}

	product, err := c.service.UpdateProduct(ctx, id, request)
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "Failed to update product",
			Reference: "ERROR-3",
			Errors:    map[string]string{"error": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[models.Product]{Item: product}, http.StatusOK)
}

// @Summary		Delete product
// @Description	API untuk menghapus produk
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
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "Failed to delete product",
			Reference: "ERROR-3",
			Errors:    map[string]string{"error": err.Error()},
		}, http.StatusNotFound)
		return
	}

	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{}, http.StatusOK)
}
