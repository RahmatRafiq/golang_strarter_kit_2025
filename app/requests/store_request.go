package requests

type StoreRequestPut struct {
	ID      uint   `json:"id" form:"id"`
	Name    string `json:"name" form:"name" example:"John Doe"`
	Phone   string `json:"phone" form:"phone" example:"08123456789"`
	Address string `json:"address" form:"address" example:"Jl. Raya No. 1"`
	City    string `json:"city" form:"city" example:"Jakarta"`
	State   string `json:"state" form:"state" example:"DKI Jakarta"`
	Country string `json:"country" form:"country" example:"Indonesia"`
	Zip     string `json:"zip" form:"zip" example:"12345"`
}

type StoreRequestMember struct {
	MemberID uint `json:"member_id" form:"member_id" binding:"required" validate:"required" example:"1"`
}