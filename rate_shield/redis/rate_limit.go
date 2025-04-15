package redisClient

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRateLimit struct {
	client *redis.ClusterClient
}

func (r RedisRateLimit) JSONSet(key string, val interface{}) error {
	return r.client.JSONSet(ctx, key, ".", val).Err()
}

func (r RedisRateLimit) JSONGet(key string) (string, bool, error) {
	res, err := r.client.JSONGet(ctx, key, ".").Result()
	if err == redis.Nil || len(res) == 0 {
		return "", false, nil
	} else if err != nil {
		return "", false, err
	}

	return res, true, nil
}

func (r RedisRateLimit) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

func (r RedisRateLimit) Delete(key string) error {
	return r.client.Del(ctx, key).Err()
}
