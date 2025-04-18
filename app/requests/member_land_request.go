package requests

type MemberLandRequestPut struct {
	ID          uint     `form:"id" json:"id"`
	MemberID    uint     `form:"member_id" json:"member_id"`
	JenisTanam  []string `form:"jenis_tanam" json:"jenis_tanam" type:"array:string"`
	JumlahPanen int      `form:"jumlah_panen" json:"jumlah_panen"`
	LuasLahan   float32  `form:"luas_lahan" json:"luas_lahan"`
	Alamat      string   `form:"alamat" json:"alamat"`
	Latitude    float32  `form:"latitude" json:"latitude"`
	Longitude   float32  `form:"longitude" json:"longitude"`
	Sertifikat  string   `form:"sertifikat" json:"sertifikat" swaggertype:"string" format:"base64"`
	Description string   `form:"description" json:"description"`
}
