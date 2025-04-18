package requests

import "time"

type ProductRequest struct {
	ID          uint      `json:"id,omitempty" form:"id,omitempty"`
	Reference   string    `json:"reference" form:"reference" example:"PRD001"`
	StoreID     uint      `json:"store_id" form:"store_id" binding:"required" example:"1"`
	CategoryID  uint      `json:"category_id" form:"category_id" binding:"required" example:"2"`
	Name        string    `json:"name" form:"name" binding:"required" example:"Product Name"`
	Description string    `json:"description" form:"description" example:"A brief description of the product"`
	Price       float64   `json:"price" form:"price" binding:"required" example:"99.99"`
	Margin      float64   `json:"margin" form:"margin" example:"10.0"`
	Stock       int       `json:"stock" form:"stock" example:"100"`
	Sold        int       `json:"sold" form:"sold" example:"10"`
	Images      []string  `json:"images" form:"images" type:"array:string"`
	ReceivedAt  time.Time `json:"received_at" form:"received_at" example:"2023-10-10T00:00:00Z"`
}
