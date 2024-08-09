package limiter

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/x-sushant-x/RateShield/models"
	"github.com/x-sushant-x/RateShield/redis"
)

var (
	TokenBucketI *TokenBuckets
	once         sync.Once
)

type TokenBuckets struct{}

func (b *TokenBuckets) SpawnNew(ip string, endpoint string, capacity int) error {
	bucket := &models.Bucket{
		ClientIP:        ip,
		CreatedAt:       time.Now().Unix(),
		Capacity:        capacity,
		AvailableTokens: capacity,
		Endpoint:        endpoint,
		TokenAddRate:    2,
		TokenAddTime:    30,
	}

	key := ip + ":" + endpoint
	err := redis.SetJSONObject(key, bucket)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (b *TokenBuckets) GetBucket(ip, endpoint string) (models.Bucket, error) {
	key := ip + ":" + endpoint
	bytes, err := redis.GetJSONObject(key)
	if err != nil {

		fmt.Println(err)
		return models.Bucket{}, err
	}

	var bucket models.Bucket
	err = json.Unmarshal(bytes, &bucket)
	if err != nil {
		fmt.Println(err)
		return models.Bucket{}, err
	}

	return bucket, nil
}

func (b *TokenBuckets) AddTokens() {}

func (b *TokenBuckets) ProcessRequest(ip, endpoint string) bool {
	bucket, err := b.GetBucket(ip, endpoint)
	if err != nil {
		fmt.Println(err)
		return false
	}

	isTokenAvailable := b.checkAvailiblity(bucket)
	if !isTokenAvailable {
		return false
	}

	b.removeToken(bucket)
	return true
}

func (b *TokenBuckets) checkAvailiblity(bucket models.Bucket) bool {
	return bucket.AvailableTokens > 0
}

func (b *TokenBuckets) removeToken(bucket models.Bucket) {
	bucket.AvailableTokens = bucket.AvailableTokens - 1

}
