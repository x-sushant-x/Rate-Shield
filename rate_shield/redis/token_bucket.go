package redisClient

import (
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
)

type RedisTokenBucket struct {
	client *redis.Client
}

func (r RedisTokenBucket) JSONSet(key string, val interface{}) error {
	return r.client.JSONSet(ctx, key, ".", val).Err()
}

func (r RedisTokenBucket) JSONGet(key string) (*models.Bucket, bool, error) {
	res, err := r.client.JSONGet(ctx, key, ".").Result()
	if err == redis.Nil || len(res) == 0 {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	var bucket models.Bucket
	err = json.Unmarshal([]byte(res), &bucket)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling bucket from Redis")
		return nil, true, err
	}

	return &bucket, true, nil
}

func (r RedisTokenBucket) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}
