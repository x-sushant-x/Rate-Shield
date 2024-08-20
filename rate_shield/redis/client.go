package redisClient

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
)

var (
	TokenBucketClient *redis.Client
	RuleClient        *redis.Client
	ctx               = context.Background()
)

func Connect() error {
	c := redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
	)

	TokenBucketClient = c

	bucketCmd := TokenBucketClient.Ping(ctx)
	if bucketCmd.Err() != nil {
		log.Error().Msgf("unable to connect to redis (token bucket db): " + bucketCmd.Err().Error())
		return bucketCmd.Err()
	} else {
		log.Info().Msg("Connected to redis successfully (token bucket db)")
	}

	ruleDB := redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       1,
		},
	)

	RuleClient = ruleDB

	ruleCmd := RuleClient.Ping(ctx)
	if ruleCmd.Err() != nil {
		log.Error().Msgf("unable to connect to redis (rules db): " + bucketCmd.Err().Error())
		return ruleCmd.Err()
	} else {
		log.Info().Msg("Connected to redis successfully (rules db)")
	}
	return nil
}

func SetJSONObject(key string, val interface{}) error {
	err := TokenBucketClient.JSONSet(ctx, key, ".", val).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetJSONObject(key string) ([]byte, error) {
	res, err := TokenBucketClient.JSONGet(ctx, key, ".").Result()
	if err != nil {
		return nil, err
	}
	return []byte(res), nil
}

func Get(key string) ([]byte, bool, error) {
	res, err := TokenBucketClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	return []byte(res), true, nil
}

func GetRule(key string) (models.Rule, bool, error) {
	res, err := RuleClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return models.Rule{}, false, nil
	} else if err != nil {
		return models.Rule{}, false, err
	}

	var rule models.Rule
	err = json.Unmarshal([]byte(res), &rule)
	if err != nil {
		return models.Rule{}, false, nil
	}

	return rule, true, nil
}

func GetAllRuleKeys() ([]string, bool, error) {
	res, err := RuleClient.Keys(ctx, "*").Result()
	if err != redis.Nil {
		return nil, false, nil
	}

	return res, true, nil
}

func SetRule(key string, val interface{}) error {
	err := RuleClient.JSONSet(ctx, key, ".", val).Err()
	if err != nil {
		return err
	}
	return nil
}

func DeleteRule(key string) error {
	err := RuleClient.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
