package models

type FixedWindowCounter struct {
	Endpoint       string `json:"endpoint"`
	ClientIP       string `json:"client_ip"`
	CreatedAt      int64  `json:"created_at"`
	MaxRequests    int64  `json:"max_requests"`
	CurrRequests   int64  `json:"current_requests"`
	Window         int    `json:"window"`
	LastAccessTime int64  `json:"last_access_time"`
}
