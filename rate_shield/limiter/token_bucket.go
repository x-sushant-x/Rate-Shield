package limiter

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
	redisClient "github.com/x-sushant-x/RateShield/redis"
	"github.com/x-sushant-x/RateShield/service"
	"github.com/x-sushant-x/RateShield/utils"
)

const (
	BucketExpireTime = time.Second * 60
)

type TokenBucketService struct {
	redisClient          redisClient.RedisRateLimiterClient
	errorNotificationSVC service.ErrorNotificationSVC
}

func NewTokenBucketService(client redisClient.RedisRateLimiterClient, errorNotificationSVC service.ErrorNotificationSVC) TokenBucketService {
	return TokenBucketService{
		redisClient:          client,
		errorNotificationSVC: errorNotificationSVC,
	}
}

func (t *TokenBucketService) addTokensToBucket(key string) {
	bucket, found, err := t.getBucket(key)
	if err != nil {
		t.sendGetBucketErrorNotification(key, err)
		return
	}

	if !found {
		return
	}

	if bucket.AvailableTokens < bucket.Capacity {
		tokensToAdd := bucket.Capacity - bucket.AvailableTokens
		bucket.AvailableTokens += min(bucket.TokenAddRate, tokensToAdd)

		if err := t.redisClient.JSONSet(key, bucket); err != nil {
			t.sendSetBucketErrorNotification(key, bucket, err)
		}
	}
}

func (t *TokenBucketService) createBucket(ip, endpoint string, capacity, tokenAddRate int) (*models.Bucket, error) {
	if err := utils.ValidateCreateBucketReq(ip, endpoint, capacity, tokenAddRate); err != nil {
		return nil, err
	}

	b := &models.Bucket{
		ClientIP:        ip,
		CreatedAt:       time.Now().Unix(),
		Capacity:        capacity,
		AvailableTokens: capacity,
		Endpoint:        endpoint,
		TokenAddRate:    tokenAddRate,
	}

	err := t.saveBucket(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (t *TokenBucketService) createBucketFromRule(ip, endpoint string, rule *models.Rule) (*models.Bucket, error) {
	b, err := t.createBucket(ip, endpoint, int(rule.TokenBucketRule.BucketCapacity), int(rule.TokenBucketRule.TokenAddRate))
	if err != nil {
		return nil, err
	}
	return b, nil
}

func parseKey(key string) (string, string) {
	parts := strings.Split(key, ":")
	return parts[0], parts[1]
}

func (t *TokenBucketService) spawnNewBucket(key string, rule *models.Rule) (*models.Bucket, error) {
	ip, endpoint := parseKey(key)
	return t.createBucketFromRule(ip, endpoint, rule)
}

func (t *TokenBucketService) getBucket(key string) (*models.Bucket, bool, error) {
	key = "token_bucket_" + key
	data, found, err := t.redisClient.JSONGet(key)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching bucket from Redis")
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	tokenBucket, err := utils.Unmarshal[models.Bucket]([]byte(data))
	if err != nil {
		log.Err(err).Msg(err.Error())
	}

	return &tokenBucket, true, nil
}

func (t *TokenBucketService) addTokens() {
	ctx := context.TODO()
	keys, err := redisClient.TokenBucketClient.Keys(ctx, "*").Result()
	if err != nil {
		log.Error().Err(err).Msg("Unable to get Redis keys")
		return
	}

	for _, key := range keys {
		t.addTokensToBucket(key)
	}
}

func (t *TokenBucketService) processRequest(key string, rule *models.Rule) *models.RateLimitResponse {
	bucket, found, err := t.getBucket(key)
	if err != nil {
		log.Error().Msgf("error while getting bucket %s" + err.Error())
		return utils.BuildRateLimitErrorResponse(500)
	}

	if !found {
		b, err := t.spawnNewBucket(key, rule)
		if err != nil {
			return utils.BuildRateLimitErrorResponse(500)
		}
		bucket = b
	}

	if bucket.AvailableTokens <= 0 {
		return utils.BuildRateLimitErrorResponse(429)
	}

	bucket.AvailableTokens--

	if err := t.saveBucket(bucket); err != nil {
		return utils.BuildRateLimitErrorResponse(500)
	}

	return &models.RateLimitResponse{
		RateLimit_Limit:     int64(bucket.Capacity),
		RateLimit_Remaining: int64(bucket.AvailableTokens),
		Success:             true,
		HTTPStatusCode:      http.StatusOK,
	}
}

func (t *TokenBucketService) saveBucket(bucket *models.Bucket) error {
	key := "token_bucket_" + bucket.ClientIP + ":" + bucket.Endpoint
	if err := t.redisClient.JSONSet(key, bucket); err != nil {
		log.Error().Err(err).Msg("Error saving new bucket to Redis")
		return err
	}

	if err := t.redisClient.Expire(key, BucketExpireTime); err != nil {
		log.Error().Err(err).Msg("Error setting bucket expiration in Redis")
		return err
	}

	return nil
}

func (t *TokenBucketService) startAddTokenJob() {
	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	_, err = s.NewJob(gocron.DurationJob(TokenAddTime), gocron.NewTask(func() {
		t.addTokens()
	}))

	if err != nil {
		panic(err)
	}

	s.Start()
}

func (t *TokenBucketService) sendGetBucketErrorNotification(key string, err error) {
	customError := fmt.Sprintf("Unable to get bucket with key: %s got error: %s", key, err.Error())
	t.errorNotificationSVC.SendErrorNotification(customError, time.Now(), "Nil", "Nil", models.Rule{})
	log.Error().Err(err).Msg("Error fetching bucket")
}

func (t *TokenBucketService) sendSetBucketErrorNotification(key string, bucket *models.Bucket, err error) {
	customError := fmt.Sprintf("Unable save bucket with key: %s and data: %+v got error: %s", key, bucket, err.Error())
	t.errorNotificationSVC.SendErrorNotification(customError, time.Now(), "Nil", "Nil", models.Rule{})
	log.Error().Err(err).Msg("Error saving updated bucket to Redis")
}
