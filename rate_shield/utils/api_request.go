package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func ExtractAuthToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")

	if header == "" || strings.HasPrefix(header, "Bearer ") {
		return "", fmt.Errorf("user not logged in")
	}

	return header[7:], nil
}
