package requests

type ProductOfferingRequestPut struct {
	Name        string  `json:"name" form:"name" binding:"required" example:"Product Name" validate:"required"`
	Description string  `json:"description" form:"description" binding:"required" example:"Product Description" validate:"required"`
	Status      string  `json:"status" form:"status" binding:"required" example:"active" validate:"required" enums:"active,inactive"`
	ID          uint    `json:"id" form:"id"`
	Price       float32 `json:"price" form:"price" binding:"required" example:"100000" validate:"required"`
}
