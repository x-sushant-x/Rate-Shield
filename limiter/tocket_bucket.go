package limiter

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/config"
	"github.com/x-sushant-x/RateShield/models"
	"github.com/x-sushant-x/RateShield/redis"
)

var (
	TokenBucketI *TokenBuckets
)

type TokenBuckets struct{}

func (b *TokenBuckets) SpawnNew(key string, capacity int) (models.Bucket, error) {
	ip := strings.Split(key, ":")[0]
	endpoint := strings.Split(key, ":")[1]

	bucket := &models.Bucket{
		ClientIP:        ip,
		CreatedAt:       time.Now().Unix(),
		Capacity:        capacity,
		AvailableTokens: capacity,
		Endpoint:        endpoint,
		TokenAddRate:    2,
		TokenAddTime:    60,
	}

	err := redis.SetJSONObject(key, bucket)
	if err != nil {
		fmt.Println(err)
		return models.Bucket{}, err
	}

	return *bucket, nil
}

func (b *TokenBuckets) GetBucket(key string) (models.Bucket, error) {
	bytes, err := redis.GetJSONObject(key)
	if err != nil {
		fmt.Println(err)
		return models.Bucket{}, err
	}

	// No bucket is available in redis
	if len(bytes) == 0 {
		bucket, err := b.SpawnNew(key, config.Config.TokenBucketCapacity)
		if err != nil {
			log.Error().Msgf("error spawning new bucket: %s" + err.Error())
			return models.Bucket{}, err
		}
		return bucket, nil
	}

	var bucket models.Bucket
	err = json.Unmarshal(bytes, &bucket)
	if err != nil {
		log.Error().Msgf("error while marshalling redis bucket: %s" + err.Error())
		return models.Bucket{}, err
	}

	return bucket, nil
}

func (b *TokenBuckets) AddTokens() {
	ctx := context.TODO()
	keys, err := redis.Client.Keys(ctx, "*").Result()
	if err != nil {
		log.Error().Msgf("unable to get all redis keys: %s", err.Error())
		return
	}

	if len(keys) == 0 {
		log.Error().Msgf("no buckets found in redis")
		return
	}

	for _, key := range keys {
		bytes, err := redis.GetJSONObject(key)
		if err != nil {
			log.Error().Msgf("error fetching bucket: %s", err.Error())
			continue
		}

		var bucket models.Bucket
		err = json.Unmarshal(bytes, &bucket)
		if err != nil {
			log.Error().Msgf("error marshalling bucket: %s", err.Error())
			continue
		}

		if bucket.AvailableTokens < bucket.Capacity {
			if bucket.AvailableTokens+bucket.TokenAddRate > bucket.Capacity {
				bucket.AvailableTokens = bucket.Capacity
			} else {
				bucket.AvailableTokens += bucket.TokenAddRate
			}
			err := redis.SetJSONObject(key, &bucket)
			if err != nil {
				log.Error().Msgf("unable to save updated bucket in redis: %s", err.Error())
				continue
			}
		}
	}
}

func (b *TokenBuckets) ProcessRequest(ip, endpoint string) bool {
	key := ip + ":" + endpoint
	bucket, err := b.GetBucket(key)
	if err != nil {
		log.Error().Msgf("error while getting bucket %s" + err.Error())
		return false
	}

	isTokenAvailable := b.checkAvailiblity(bucket)
	if !isTokenAvailable {
		return false
	}

	bucket.AvailableTokens = bucket.AvailableTokens - 1
	err = redis.SetJSONObject(key, bucket)
	if err != nil {
		log.Error().Msgf("error while storing removed token bucket to redis: %s" + err.Error())
		return false
	}
	return true
}

func (b *TokenBuckets) checkAvailiblity(bucket models.Bucket) bool {
	return bucket.AvailableTokens > 0
}
