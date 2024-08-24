package api

import (
	"net/http"

	"github.com/x-sushant-x/RateShield/limiter"
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

	code := h.limiterSvc.CheckLimit(ip, endpoint)

	switch code {
	case 200:
		w.WriteHeader(http.StatusOK)
	case 429:
		w.WriteHeader(http.StatusTooManyRequests)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
