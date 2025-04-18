package responses

type CheckNIK struct {
	IsRegistered bool `json:"is_registered"`
	IsRegisteredOnStore bool `json:"is_registered_on_store"`
}