package utils

import (
	"encoding/json"
	"net/http"
)

func InternalError(w http.ResponseWriter) {
	msg := map[string]string{
		"status": "fail",
		"error":  "Internal Server Error",
	}

	w.WriteHeader(http.StatusInternalServerError)
	bytes, _ := json.Marshal(msg)
	w.Write(bytes)
}

func BadRequestError(w http.ResponseWriter) {
	msg := map[string]string{
		"status": "fail",
		"error":  "Invalid Request Body",
	}

	w.WriteHeader(http.StatusInternalServerError)
	bytes, _ := json.Marshal(msg)
	w.Write(bytes)
}

func SuccessResponse(data interface{}, w http.ResponseWriter) {
	msg := map[string]interface{}{
		"status": "success",
		"data":   data,
	}

	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(msg)
	w.Write(bytes)
}
