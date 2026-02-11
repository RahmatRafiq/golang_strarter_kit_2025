package requests

type CostBudgetPlanRequestPut struct {
	Status       string `json:"status" form:"status" binding:"required" validate:"required" example:"draft"`
	ApprovedAt   string `json:"approved_at" example:"2021-01-01 00:00:00"`
	JenisTanam   string `json:"jenis_tanam" form:"jenis_tanam" binding:"required" validate:"required" example:"Jenis Tanam"`
	ID           uint   `json:"id" form:"id"`
	MemberID     uint   `json:"member_id" form:"member_id" binding:"required" validate:"required" example:"1"`
	StoreID      uint   `json:"store_id" form:"store_id" binding:"required" validate:"required" example:"1"`
	MemberLandID uint   `json:"member_land_id" form:"member_land_id" binding:"required" validate:"required" example:"1"`
}

type CostBudgetPlanWizardRequestPut struct {
}

type CostBudgetPlanItemRequestPut struct {
	Description       string  `json:"description" form:"description" binding:"required" validate:"required" example:"Description"`
	ID                uint    `json:"id" form:"id"`
	ProductOfferingID uint    `json:"product_offering_id" form:"product_offering_id" binding:"required" validate:"required" example:"1"`
	PhaseID           uint    `json:"phase_id" form:"phase_id" binding:"required" validate:"required" example:"1"`
	Quantity          float32 `json:"quantity" form:"quantity" binding:"required" validate:"required" example:"1"`
}

type CostBudgetPlanDetailHistoryRequest struct {
	Description string  `json:"description" form:"description" binding:"required" validate:"required"`
	ID          uint    `json:"id" form:"id"`
	ProductID   uint    `json:"product_id" form:"product_id" binding:"required" validate:"required"`
	Quantity    float32 `json:"quantity" form:"quantity" binding:"required" validate:"required"`
}
