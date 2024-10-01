package redisClient

import (
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
)

type RedisFixedWindow struct {
	client *redis.Client
}

func (r RedisFixedWindow) Expire(key string, expiration time.Duration) error {
	_, err := r.client.Expire(ctx, key, expiration).Result()
	return err
}

func (r RedisFixedWindow) JSONGet(key string) (*models.FixedWindowCounter, bool, error) {
	res, err := r.client.JSONGet(ctx, key, ".").Result()
	if err == redis.Nil || len(res) == 0 {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	fw := models.FixedWindowCounter{}

	err = json.Unmarshal([]byte(res), &fw)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling fixed window from Redis")
		return nil, true, err
	}

	return &fw, true, nil
}

func (r RedisFixedWindow) JSONSet(key string, val interface{}) error {
	err := r.client.JSONSet(ctx, key, ".", val).Err()
	if err != nil {
		return err
	}
	return nil
}
