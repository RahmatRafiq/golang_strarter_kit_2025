package helpers

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func GenerateReference(code string) string {
	ref, err := uuid.NewV7()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate reference UUID")
	}
	return code + "-" + ref.String()
}
