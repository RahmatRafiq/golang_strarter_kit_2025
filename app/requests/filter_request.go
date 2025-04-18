package requests

type FilterRequest struct {
	Search         *string `form:"search" json:"search"`
	OrderBy        *string `form:"order_by" json:"order_by"`
	OrderDirection *string `form:"order_direction" json:"order_direction" enums:"asc,desc"`
	Page           *int    `form:"page" json:"page"`
	Limit          *int    `form:"limit" json:"limit"`
	Offset         *int    `form:"offset" json:"offset"`
}
