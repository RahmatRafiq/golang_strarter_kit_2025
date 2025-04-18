package requests

type MemberRequestPut struct {
	ID           uint   `json:"id" form:"id"`
	Nama         string `json:"nama" form:"nama" binding:"required" example:"John Doe" validate:"required"`
	Phone        string `json:"phone" form:"phone" binding:"required" example:"08123456789" validate:"required"`
	Alamat       string `json:"alamat" form:"alamat" binding:"required" example:"Jl. Raya No. 1" validate:"required"`
	NIK          string `json:"nik" form:"nik" binding:"required" example:"1234567890123456" validate:"required"`
	JenisKelamin string `json:"jenis_kelamin" form:"jenis_kelamin" binding:"required" example:"pria" validate:"required"`
	FotoKTP      string `json:"foto_ktp" form:"foto_ktp" binding:"required" example:"base64" validate:"required"`
}

type MemberRequestPutLands struct {
	MemberLands []MemberLandRequestPut `json:"member_lands" form:"member_lands" binding:"required" validate:"required"`
}

type MemberRequestCheckNIK struct {
	NIK     string `json:"nik" form:"nik" binding:"required" example:"1234567890123456" validate:"required"`
	StoreID uint   `json:"store_id" form:"store_id" binding:"required" example:"1" validate:"required"`
}
