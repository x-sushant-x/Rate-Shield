package limiter

import (
	"time"

	"github.com/x-sushant-x/ThrottleWatch/config"
)

type TokenBuckets struct {
	buckets map[string]*Bucket
}

type Bucket struct {
	ClientIP  string
	CreatedAt int64
	Capacity  int
	Available int
}

func (b *TokenBuckets) SpawnNew(ip string, capacity int) *Bucket {
	bucket := &Bucket{
		ClientIP:  ip,
		CreatedAt: time.Now().Unix(),
		Capacity:  capacity,
		Available: capacity,
	}

	b.buckets[ip] = bucket
	return bucket
}

func (b *TokenBuckets) AddTokens() {
	addRate := config.Config.TokenAddingRate

	for _, bucket := range b.buckets {
		if bucket.Available < bucket.Capacity {
			remaining := bucket.Capacity - bucket.Available
			if remaining >= addRate {
				bucket.Available += addRate
				continue
			}
			bucket.Available = bucket.Capacity
		}
	}
}

func (b *TokenBuckets) ProcessRequest(ip string) bool {
	bucket := b.GetBucket(ip)

	isTokenAvailable := b.checkAvailiblity(bucket)
	if !isTokenAvailable {
		return false
	}

	b.removeToken(bucket)
	return true
}

func (b *TokenBuckets) GetBucket(ip string) *Bucket {
	bucket, found := b.buckets[ip]
	if !found {
		capacity := config.Config.TokenBucketCapacity
		return b.SpawnNew(ip, capacity)
	}
	return bucket
}

func (b *TokenBuckets) checkAvailiblity(bucket *Bucket) bool {
	return bucket.Available > 0
}

func (b *TokenBuckets) removeToken(bucket *Bucket) {
	bucket.Available = bucket.Available - 1
}
