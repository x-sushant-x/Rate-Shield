package redisClient

import (
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/x-sushant-x/RateShield/models"
)

type RedisRules struct {
	client *redis.Client
}

func (r RedisRules) DeleteRule(key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisRules) GetAllRuleKeys() ([]string, bool, error) {
	res, err := r.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, false, nil
	}

	return res, true, nil
}

func (r RedisRules) GetRule(key string) (*models.Rule, bool, error) {
	res, err := r.client.JSONGet(ctx, key).Result()
	if err == redis.Nil {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	var rule models.Rule
	err = json.Unmarshal([]byte(res), &rule)
	if err != nil {
		return nil, false, nil
	}

	return &rule, true, nil
}

func (r RedisRules) SetRule(key string, val interface{}) error {
	err := r.client.JSONSet(ctx, key, ".", val).Err()
	if err != nil {
		return err
	}
	return nil
}
