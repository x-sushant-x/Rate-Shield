package models

type Bucket struct {
	Endpoint        string `json:"endpoint"`
	Capacity        int    `json:"capacity"`
	TokenAddRate    int    `json:"token_add_rate"`
	ClientIP        string `json:"client_ip"`
	CreatedAt       int64  `json:"created_at"`
	AvailableTokens int    `json:"available_tokens"`
}

type Buckets struct {
	Buckets []Buckets `json:"buckets"`
}
