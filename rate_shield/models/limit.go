package models

type RateLimitResponse struct {
	RateLimit_Limit     int64
	RateLimit_Remaining int64
	Success             bool
	HTTPStatusCode      int
}
