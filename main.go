package main

import (
	"golang_starter_kit_2025/bootstrap"
)

func main() {
	//	@securityDefinitions.apikey	X-Api-Key
	//	@in							header
	//	@name						X-Api-Key

	//	@securityDefinitions.apikey	Bearer
	//	@in							header
	//	@name						Authorization

	// @Security X-Api-Key
	// @Security Bearer

	bootstrap.Init()
}
