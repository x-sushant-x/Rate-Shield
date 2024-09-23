package utils

import (
	"net/http"

	"github.com/x-sushant-x/RateShield/models"
)

func BuildRateLimitErrorResponse(statusCode int) *models.RateLimitResponse {
	return &models.RateLimitResponse{
		RateLimit_Limit:     -1,
		RateLimit_Remaining: -1,
		Success:             false,
		HTTPStatusCode:      statusCode,
	}
}

func BuildRateLimitSuccessResponse(limit, remaining int64) *models.RateLimitResponse {
	return &models.RateLimitResponse{
		RateLimit_Limit:     limit,
		RateLimit_Remaining: remaining,
		Success:             true,
		HTTPStatusCode:      http.StatusOK,
	}
}
