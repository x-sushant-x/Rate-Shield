package redisClient

import (
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
)

const (
	redisRuleUpdateChannel = "rules-update"
)

type RedisRules struct {
	client *redis.Client
}

// GetClient returns the underlying Redis client
func (r RedisRules) GetClient() *redis.Client {
	return r.client
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
		return nil, false, err
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

func (r RedisRules) PublishMessage(channel, msg string) error {
	return r.client.Publish(ctx, channel, msg).Err()
}

func (r RedisRules) ListenToRulesUpdate(updatesChannel chan string) {
	pubsub := r.client.Subscribe(ctx, redisRuleUpdateChannel)
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Err(err).Msg("Error while listening for rule updates")
			continue
		}

		if msg.Channel == redisRuleUpdateChannel {
			updatesChannel <- "UpdateRules"
		}
	}
}
