package redisClient

import (
	"time"

	"github.com/x-sushant-x/RateShield/models"
)

type RedisTokenBucketClient interface {
	JSONSet(key string, val interface{}) error
	JSONGet(key string) (*models.Bucket, bool, error)
	Expire(key string, expiration time.Duration) error
}

type RedisRuleClient interface {
	JSONSet(key string, path string, val interface{}) error
	JSONGet(key string, path string) (string, error)
	Keys(pattern string) ([]string, error)
	Del(key string) error
}

type RedisFixedWindowClient interface {
	JSONSet(key string, val interface{}) error
	JSONGet(key string) (*models.FixedWindowCounter, bool, error)
	Expire(key string, expireTime time.Duration) error
}
