package redisClient

import (
	"context"
	"time"

	"github.com/x-sushant-x/RateShield/models"
)

type RedisRuleClient interface {
	GetRule(key string) (*models.Rule, bool, error)
	GetAllRuleKeys() ([]string, bool, error)
	SetRule(key string, val interface{}) error
	DeleteRule(key string) error
	PublishMessage(channel, msg string) error
	ListenToRulesUpdate(udpatesChannel chan string)
}

type RedisRateLimiterClient interface {
	JSONSet(key string, val interface{}) error
	JSONGet(key string) (string, bool, error)
	Expire(key string, expireTime time.Duration) error
	Delete(key string) error
}

// SlidingWindowClient defines operations needed for sliding window rate limiting
type SlidingWindowClient interface {
	ZRemRangeByScore(ctx context.Context, key, min, max string) error
	ZCount(ctx context.Context, key, min, max string) (int64, error)
	ZAdd(ctx context.Context, key string, timestamp int64) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
}
