package api

import (
	"encoding/json"
	"net/http"

	"github.com/x-sushant-x/RateShield/service"
	"github.com/x-sushant-x/RateShield/utils"
)

type AuthAPIHandler struct {
	authSvc service.AuthService
}

func NewAuthAPIHandler(authSvc service.AuthService) AuthAPIHandler {
	return AuthAPIHandler{
		authSvc: authSvc,
	}
}

func (ah *AuthAPIHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	email := params.Get("email")
	pass := params.Get("password")

	if email == "" || pass == "" {
		utils.BadRequestError(w)
		return
	}

	err := ah.authSvc.LoginUser(email, pass)
	if err != nil {
		msg := map[string]interface{}{
			"status": "error",
			"data":   err.Error(),
		}

		w.WriteHeader(http.StatusOK)
		bytes, _ := json.Marshal(msg)
		w.Write(bytes)
	}
}

func (ah *AuthAPIHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	email := params.Get("email")
	pass := params.Get("password")

	if email == "" || pass == "" {
		utils.BadRequestError(w)
		return
	}

	err := ah.authSvc.CreateUser(email, pass)
	if err != nil {
		msg := map[string]any{
			"status": "error",
			"data":   err.Error(),
		}

		w.WriteHeader(http.StatusOK)
		bytes, _ := json.Marshal(msg)
		w.Write(bytes)
	}

	utils.SuccessResponse("Account created successfully", w)
}
