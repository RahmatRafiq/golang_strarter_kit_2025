package helpers

import (
	"log"

	"github.com/google/uuid"
)

func GenerateReference(code string) string {
	ref, err := uuid.NewV7()
	if err != nil {
		log.Fatal(err)
	}
	return code + "-" + ref.String()
}
