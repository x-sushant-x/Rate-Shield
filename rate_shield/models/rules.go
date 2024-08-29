package models

type Rule struct {
	Strategy       string `json:"strategy"`
	APIEndpoint    string `json:"endpoint"`
	BucketCapacity int64  `json:"bucket_capacity"` // For Token Bucket
	TokenAddRate   int64  `json:"token_add_rate"`  // For Token Bucket
	HTTPMethod     string `json:"http_method"`

	MaxRequests int64 `json:"max_requests"` // For Sliding Window Counter
	Window      int   `json:"window"`       // For Sliding Window Counter (in seconds)
}

type DeleteRuleDTO struct {
	RuleKey string `json:"rule_key"`
}
