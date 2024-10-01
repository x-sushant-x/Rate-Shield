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
	GetRule(key string) (*models.Rule, bool, error)
	GetAllRuleKeys() ([]string, bool, error)
	SetRule(key string, val interface{}) error
	DeleteRule(key string) error
}

type RedisFixedWindowClient interface {
	JSONSet(key string, val interface{}) error
	JSONGet(key string) (*models.FixedWindowCounter, bool, error)
	Expire(key string, expireTime time.Duration) error
}
