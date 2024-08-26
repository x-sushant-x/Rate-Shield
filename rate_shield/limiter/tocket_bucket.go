package limiter

import (
	"context"
	"encoding/json"
	"fmt"
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

func (b *TokenBucketService) spawnNewBucket(key string) (TokenBucketService, error) {
	ip, endpoint := parseKey(key)

	rule, found, err := redisClient.GetRule(endpoint)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching rule from Redis")
		return TokenBucketService{}, err
	}

	if !found {
		return createBucket(ip, endpoint, config.Config.TokenBucketCapacity, config.Config.TokenAddingRate), nil
	}

	return b.createBucketFromRule(ip, endpoint, rule), nil
}

func (b *TokenBucketService) GetBucket(key string) (TokenBucketService, error) {
	data, found, err := redisClient.GetJSONObject(key)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching bucket from Redis")

	}

	if len(data) == 0 || !found {
		return b.spawnNewBucket(key)
	}

	return b.unmarshalBucket(data)
}

func (b *TokenBucketService) AddTokens() {
	ctx := context.TODO()
	keys, err := redisClient.TokenBucketClient.Keys(ctx, "*").Result()
	if err != nil {
		log.Error().Err(err).Msg("Unable to get Redis keys")
		return
	}

	for _, key := range keys {
		b.addTokensToBucket(key)
	}
}

func (b *TokenBucketService) ProcessRequest(ip, endpoint string) int {
	key := ip + ":" + endpoint

	bucket, err := b.GetBucket(key)
	if err != nil {
		log.Error().Msgf("error while getting bucket %s" + err.Error())
		return http.StatusInternalServerError
	}

	if !b.checkAvailiblity(bucket) {
		return http.StatusTooManyRequests
	}

	bucket.AvailableTokens--

	if err := b.saveBucket(); err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func (b *TokenBucketService) checkAvailiblity(bucket TokenBucketService) bool {
	return bucket.AvailableTokens > 0
}

func parseKey(key string) (string, string) {
	parts := strings.Split(key, ":")
	return parts[0], parts[1]
}

func createBucket(ip, endpoint string, capacity, tokenAddRate int) TokenBucketService {
	b := TokenBucketService{
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

func (t *TokenBucketService) createBucketFromRule(ip, endpoint string, rule models.Rule) TokenBucketService {
	return createBucket(ip, endpoint, int(rule.BucketCapacity), int(rule.TokenAddRate))
}

func (t *TokenBucketService) saveBucket() error {
	key := t.ClientIP + ":" + t.Endpoint
	if key == ":" {
		log.Error().Msg("Bucket key is empty")
		fmt.Println(t)
	}
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

func (t *TokenBucketService) unmarshalBucket(data []byte) (TokenBucketService, error) {
	var bucket TokenBucketService
	if err := json.Unmarshal(data, &bucket); err != nil {
		log.Error().Err(err).Msg("Error unmarshalling bucket data")
		return TokenBucketService{}, err
	}
	return bucket, nil
}

func (b *TokenBucketService) addTokensToBucket(key string) {
	bucket, err := b.GetBucket(key)
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
