package limiter

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
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
	BucketExpireTime    = time.Second * 30
	DefaultTokenAddTime = 0
)

type TokenBucketService struct{}

func (b *TokenBucketService) spawnNewBucket(key string) (models.Bucket, error) {
	ip, endpoint := parseKey(key)

	rule, found, err := redisClient.GetRule(endpoint)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching rule from Redis")
		return models.Bucket{}, err
	}

	if !found {
		return b.createBucket(ip, endpoint, config.Config.TokenBucketCapacity, config.Config.TokenAddingRate), nil
	}

	return b.createBucketFromRule(ip, endpoint, string(rule)), nil
}

func (b *TokenBucketService) GetBucket(key string) (models.Bucket, error) {
	data, err := redisClient.GetJSONObject(key)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching bucket from Redis")
		return models.Bucket{}, err
	}

	if len(data) == 0 {
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

	if err := b.saveBucket(key, bucket); err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func (b *TokenBucketService) checkAvailiblity(bucket models.Bucket) bool {
	return bucket.AvailableTokens > 0
}

func parseKey(key string) (string, string) {
	parts := strings.Split(key, ":")
	return parts[0], parts[1]
}

func parseRule(rule string) (int, int) {
	parts := strings.Split(rule, ":")
	capacity, _ := strconv.Atoi(parts[0])
	tokenAddRate, _ := strconv.Atoi(parts[1])
	return capacity, tokenAddRate
}

func (t *TokenBucketService) createBucket(ip, endpoint string, capacity, tokenAddRate int) models.Bucket {
	bucket := models.Bucket{
		ClientIP:        ip,
		CreatedAt:       time.Now().Unix(),
		Capacity:        capacity,
		AvailableTokens: capacity,
		Endpoint:        endpoint,
		TokenAddRate:    tokenAddRate,
		TokenAddTime:    DefaultTokenAddTime,
	}

	t.saveBucket(ip, bucket)
	return bucket
}

func (t *TokenBucketService) createBucketFromRule(ip, endpoint, rule string) models.Bucket {
	capacity, tokenAddRate := parseRule(rule)
	return t.createBucket(ip, endpoint, capacity, tokenAddRate)
}

func (t *TokenBucketService) saveBucket(key string, bucket models.Bucket) error {
	if err := redisClient.SetJSONObject(key, bucket); err != nil {
		log.Error().Err(err).Msg("Error saving new bucket to Redis")
		return err
	}

	if err := redisClient.TokenBucketClient.Expire(context.Background(), key, BucketExpireTime).Err(); err != nil {
		log.Error().Err(err).Msg("Error setting bucket expiration in Redis")
		return err
	}

	return nil
}

func (t *TokenBucketService) unmarshalBucket(data []byte) (models.Bucket, error) {
	var bucket models.Bucket
	if err := json.Unmarshal(data, &bucket); err != nil {
		log.Error().Err(err).Msg("Error unmarshalling bucket data")
		return models.Bucket{}, err
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
