package limiter

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
	redisClient "github.com/x-sushant-x/RateShield/redis"
)

const (
	BucketExpireTime    = time.Second * 60
	DefaultTokenAddTime = 60
)

type TokenBucketService models.Bucket

func NewTokenBucketService() TokenBucketService {
	return TokenBucketService{}
}

func addTokensToBucket(key string) {
	bucket, found, err := getBucket(key)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching bucket")
		return
	}

	if !found {
		log.Error().Err(err).Msg("Error fetching bucket with given key")
		return
	}

	if bucket.AvailableTokens < bucket.Capacity {
		tokensToAdd := bucket.Capacity - bucket.AvailableTokens
		bucket.AvailableTokens += min(bucket.TokenAddRate, tokensToAdd)

		if err := redisClient.SetTokenBucketJSONObject(key, bucket); err != nil {
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
	return createBucket(ip, endpoint, int(rule.TokenBucketRule.BucketCapacity), int(rule.TokenBucketRule.TokenAddRate))
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

func spawnNewBucket(key string, rule models.Rule) (*TokenBucketService, error) {
	ip, endpoint := parseKey(key)
	return createBucketFromRule(ip, endpoint, rule), nil
}

func getBucket(key string) (*TokenBucketService, bool, error) {
	data, found, err := redisClient.GetTokenBucketJSONObject(key)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching bucket from Redis")
		return nil, false, err
	}

	if !found {
		log.Error().Err(err).Msg("Error bucket not found in redis")
		return nil, false, nil
	}

	bucket, err := unmarshalBucket(data)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling bucket from Redis")
		return nil, false, err
	}
	return bucket, true, nil
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

func (t *TokenBucketService) processRequest(key string, rule models.Rule) int {
	bucket, found, err := getBucket(key)
	if err != nil {
		log.Error().Msgf("error while getting bucket %s" + err.Error())
		return http.StatusInternalServerError
	}

	if !found {
		b, err := spawnNewBucket(key, rule)
		if err != nil {
			return http.StatusInternalServerError
		}
		bucket = b
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
	if err := redisClient.SetTokenBucketJSONObject(key, t); err != nil {
		log.Error().Err(err).Msg("Error saving new bucket to Redis")
		return err
	}

	if err := redisClient.TokenBucketClient.Expire(context.Background(), key, BucketExpireTime).Err(); err != nil {
		log.Error().Err(err).Msg("Error setting bucket expiration in Redis")
		return err
	}

	return nil
}
