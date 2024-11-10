package models

type Rule struct {
	Strategy                 string                    `json:"strategy"`
	APIEndpoint              string                    `json:"endpoint"`
	HTTPMethod               string                    `json:"http_method"`
	AllowOnError             bool                      `json:"allow_on_error"`
	TokenBucketRule          *TokenBucketRule          `json:"token_bucket_rule,omitempty"`
	FixedWindowCounterRule   *FixedWindowCounterRule   `json:"fixed_window_counter_rule,omitempty"`
	SlidingWindowCounterRule *SlidingWindowCounterRule `json:"sliding_window_counter_rule,omitempty"`
}

type TokenBucketRule struct {
	BucketCapacity int64 `json:"bucket_capacity"`
	TokenAddRate   int64 `json:"token_add_rate"`
}

type FixedWindowCounterRule struct {
	MaxRequests int64 `json:"max_requests"`
	Window      int   `json:"window"`
}

type DeleteRuleDTO struct {
	RuleKey string `json:"rule_key"`
}

type PaginatedRules struct {
	PageNumber  int    `json:"page_number"`
	TotalItems  int    `json:"total_items"`
	HasNextPage bool   `json:"has_next_page"`
	Rules       []Rule `json:"rules"`
}

type SlidingWindowCounterRule struct {
	MaxRequests int64 `json:"max_requests"`
	WindowSize  int   `json:"window"`
}
