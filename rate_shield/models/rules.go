package models

type Rule struct {
	Type           string `json:"type"`
	APIEndpoint    string `json:"endpoint"`
	BucketCapacity int64  `json:"bucket_capacity"`
	TokenAddRate   int64  `json:"token_add_rate"`
	HTTPMethod     string `json:"http_method"`
}

type DeleteRuleDTO struct {
	RuleKey string `json:"rule_key"`
}
