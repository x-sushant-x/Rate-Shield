package api

import (
	"fmt"
	"net/http"

	"github.com/x-sushant-x/RateShield/limiter"
	"github.com/x-sushant-x/RateShield/utils"
)

type RateLimitHandler struct {
	limiterSvc limiter.Limiter
}

func NewRateLimitHandler(limiterSvc limiter.Limiter) RateLimitHandler {
	return RateLimitHandler{
		limiterSvc: limiterSvc,
	}
}

func (h RateLimitHandler) CheckRateLimit(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("ip")
	endpoint := r.Header.Get("endpoint")

	badRequest := utils.ValidateLimitRequest(ip, endpoint)
	if badRequest != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	resp := h.limiterSvc.CheckLimit(ip, endpoint)

	switch resp.HTTPStatusCode {
	case 200:
		w.Header().Set("rate-limit", fmt.Sprint(resp.RateLimit_Limit))
		w.Header().Set("rate-limit-remaining", fmt.Sprint(resp.RateLimit_Remaining))
		w.WriteHeader(http.StatusOK)
	case 429:
		w.WriteHeader(http.StatusTooManyRequests)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
