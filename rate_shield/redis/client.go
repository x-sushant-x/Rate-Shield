package redisClient

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
)

var (
	TokenBucketClient        *redis.Client
	RuleClient               *redis.Client
	FixedWindowCounterClient *redis.Client
	ctx                      = context.Background()
)

func createNewRedisConnection(addr string, db int) (*redis.Client, error) {
	conn := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       db,
	})

	pingCmd := conn.Ping(ctx)
	if pingCmd.Err() != nil {
		return nil, pingCmd.Err()
	}
	return conn, nil
}

func Connect() error {
	ruleClient, err := createNewRedisConnection("localhost:6379", 0)
	checkError(err)
	RuleClient = ruleClient

	tokenBucketClient, err := createNewRedisConnection("localhost:6379", 1)
	checkError(err)
	TokenBucketClient = tokenBucketClient

	fixedWindowClient, err := createNewRedisConnection("localhost:6379", 2)
	checkError(err)
	FixedWindowCounterClient = fixedWindowClient

	log.Info().Msg("Connected To Redis")
	return nil
}

func SetTokenBucketJSONObject(key string, val interface{}) error {
	err := TokenBucketClient.JSONSet(ctx, key, ".", val).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetTokenBucketJSONObject(key string) ([]byte, bool, error) {
	res, err := TokenBucketClient.JSONGet(ctx, key, ".").Result()
	if err == redis.Nil || len(res) == 0 {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	return []byte(res), true, nil
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

func GetRule(key string) (*models.Rule, bool, error) {
	res, err := RuleClient.JSONGet(ctx, key).Result()
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

func GetAllRuleKeys() ([]string, bool, error) {
	res, err := RuleClient.Keys(ctx, "*").Result()
	if err != nil {
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

func checkError(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to redis")
	}
}

func SetFixedWindowJSONObject(key string, val interface{}) error {
	err := FixedWindowCounterClient.JSONSet(ctx, key, ".", val).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetFixedWindowJSONObject(key string) ([]byte, bool, error) {
	res, err := FixedWindowCounterClient.JSONGet(ctx, key, ".").Result()
	if err == redis.Nil || len(res) == 0 {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	return []byte(res), true, nil
}

func SetFixedWindowExpireTime(key string, expireTime time.Duration) error {
	_, err := FixedWindowCounterClient.Expire(ctx, key, expireTime).Result()
	return err
}
