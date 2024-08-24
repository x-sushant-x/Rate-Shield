package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseAPIBody[T any](r *http.Request) (T, error) {
	var req T
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return req, err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		return req, err
	}
	return req, nil
}
