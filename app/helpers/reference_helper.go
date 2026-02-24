package helpers

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateReference(code string) (string, error) {
	ref, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate reference UUID: %w", err)
	}
	return code + "-" + ref.String(), nil
}
