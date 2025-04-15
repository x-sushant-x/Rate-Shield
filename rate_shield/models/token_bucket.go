package models

import "time"

type Bucket struct {
	Endpoint        string    `json:"endpoint"`
	Capacity        int       `json:"capacity"`
	TokenAddRate    int       `json:"token_add_rate"`
	ClientIP        string    `json:"client_ip"`
	CreatedAt       int64     `json:"created_at"`
	AvailableTokens int       `json:"available_tokens"`
	LastRefill      time.Time `json:"last_refill"`
	RetentionTime   int16     `json:"retention_time"` // Amount of time to keep inactive bucket in redis (in seconds)
}

type Buckets struct {
	Buckets []Buckets `json:"buckets"`
}
