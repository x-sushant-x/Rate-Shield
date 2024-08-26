package limiter

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/config"
	"github.com/x-sushant-x/RateShield/models"
	redisClient "github.com/x-sushant-x/RateShield/redis"
)

var (
	TokenBucketManager    *TokenBucketService
	DefaultTokenAddRate   = config.Config.TokenAddingRate
	DefaultBucketCapacity = config.Config.TokenBucketCapacity
)

const (
	BucketExpireTime    = time.Second * 60
	DefaultTokenAddTime = 0
)

type TokenBucketService models.Bucket

func NewTokenBucketService() TokenBucketService {
	return TokenBucketService{}
}

func addTokensToBucket(key string) {
	bucket, err := getBucket(key)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching bucket")
		return
	}

	if bucket.AvailableTokens < bucket.Capacity {
		tokensToAdd := bucket.Capacity - bucket.AvailableTokens
		bucket.AvailableTokens += min(bucket.TokenAddRate, tokensToAdd)

		if err := redisClient.SetJSONObject(key, bucket); err != nil {
			log.Error().Err(err).Msg("Error saving updated bucket to Redis")
		}
	}
}

func createBucket(ip, endpoint string, capacity, tokenAddRate int) *TokenBucketService {
	b := &TokenBucketService{
		ClientIP:        ip,
		CreatedAt:       time.Now().Unix(),
		Capacity:        capacity,
		AvailableTokens: capacity,
		Endpoint:        endpoint,
		TokenAddRate:    tokenAddRate,
		TokenAddTime:    DefaultTokenAddTime,
	}

	b.saveBucket()
	return b
}

func createBucketFromRule(ip, endpoint string, rule models.Rule) *TokenBucketService {
	return createBucket(ip, endpoint, int(rule.BucketCapacity), int(rule.TokenAddRate))
}

func parseKey(key string) (string, string) {
	parts := strings.Split(key, ":")
	return parts[0], parts[1]
}

func unmarshalBucket(data []byte) (*TokenBucketService, error) {
	bucket := new(TokenBucketService)
	if err := json.Unmarshal(data, &bucket); err != nil {
		log.Error().Err(err).Msg("Error unmarshalling bucket data")
		return bucket, err
	}
	return bucket, nil
}

func spawnNewBucket(key string) (*TokenBucketService, error) {
	ip, endpoint := parseKey(key)

	rule, found, err := redisClient.GetRule(endpoint)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching rule from Redis")
		return new(TokenBucketService), err
	}

	if !found {
		return createBucket(ip, endpoint, config.Config.TokenBucketCapacity, config.Config.TokenAddingRate), nil
	}

	return createBucketFromRule(ip, endpoint, rule), nil
}

func getBucket(key string) (*TokenBucketService, error) {
	data, found, err := redisClient.GetJSONObject(key)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching bucket from Redis")

	}

	if len(data) == 0 || !found {
		return spawnNewBucket(key)
	}

	return unmarshalBucket(data)
}

func (t *TokenBucketService) addTokens() {
	ctx := context.TODO()
	keys, err := redisClient.TokenBucketClient.Keys(ctx, "*").Result()
	if err != nil {
		log.Error().Err(err).Msg("Unable to get Redis keys")
		return
	}

	for _, key := range keys {
		addTokensToBucket(key)
	}
}

func (t *TokenBucketService) processRequest(ip, endpoint string) int {
	key := ip + ":" + endpoint

	bucket, err := getBucket(key)
	if err != nil {
		log.Error().Msgf("error while getting bucket %s" + err.Error())
		return http.StatusInternalServerError
	}

	if !bucket.checkAvailiblity() {
		return http.StatusTooManyRequests
	}

	bucket.AvailableTokens--

	if err := bucket.saveBucket(); err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func (t *TokenBucketService) checkAvailiblity() bool {
	return t.AvailableTokens > 0
}

func (t *TokenBucketService) saveBucket() error {
	key := t.ClientIP + ":" + t.Endpoint
	if err := redisClient.SetJSONObject(key, t); err != nil {
		log.Error().Err(err).Msg("Error saving new bucket to Redis")
		return err
	}

	if err := redisClient.TokenBucketClient.Expire(context.Background(), key, BucketExpireTime).Err(); err != nil {
		log.Error().Err(err).Msg("Error setting bucket expiration in Redis")
		return err
	}

	return nil
}
