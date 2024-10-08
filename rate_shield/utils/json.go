package utils

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

func MarshalJSON(data any) ([]byte, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Err(err).Msgf("unable to marshal %v", data)
	}

	return dataBytes, err
}
