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

	return dataBytes, nil
}

func Unmarshal[T any](data []byte) (T, error) {
	var res T
	err := json.Unmarshal(data, &res)
	if err != nil {
		log.Err(err).Msgf("unable to unmarshal %s", string(data))
	}

	return res, nil
}
