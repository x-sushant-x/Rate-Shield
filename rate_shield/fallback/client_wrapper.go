package fallback

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	redisClient "github.com/x-sushant-x/RateShield/redis"
)

// FallbackRateLimitClient wraps Redis client with in-memory fallback
type FallbackRateLimitClient struct {
	redisClient      redisClient.RedisRateLimiterClient
	memoryStore      *InMemoryRateLimitStore
	isRedisAvailable bool
	mutex            sync.RWMutex
	lastWarningTime  time.Time
	warningInterval  time.Duration
}

// NewFallbackRateLimitClient creates a new fallback client
func NewFallbackRateLimitClient(redisClient redisClient.RedisRateLimiterClient, memoryStore *InMemoryRateLimitStore) *FallbackRateLimitClient {
	return &FallbackRateLimitClient{
		redisClient:      redisClient,
		memoryStore:      memoryStore,
		isRedisAvailable: true,
		warningInterval:  5 * time.Minute,
		lastWarningTime:  time.Now(),
	}
}

// JSONSet stores a value as JSON
func (f *FallbackRateLimitClient) JSONSet(key string, val interface{}) error {
	f.mutex.RLock()
	isAvailable := f.isRedisAvailable
	f.mutex.RUnlock()

	if isAvailable {
		err := f.redisClient.JSONSet(key, val)
		if err != nil {
			f.switchToFallback()
			return f.memoryStore.JSONSet(key, val)
		}
		return nil
	}

	f.logFallbackWarning()
	return f.memoryStore.JSONSet(key, val)
}

// JSONGet retrieves a value as JSON string
func (f *FallbackRateLimitClient) JSONGet(key string) (string, bool, error) {
	f.mutex.RLock()
	isAvailable := f.isRedisAvailable
	f.mutex.RUnlock()

	if isAvailable {
		result, found, err := f.redisClient.JSONGet(key)
		if err != nil {
			f.switchToFallback()
			return f.memoryStore.JSONGet(key)
		}
		return result, found, nil
	}

	f.logFallbackWarning()
	return f.memoryStore.JSONGet(key)
}

// Expire sets an expiration time on a key
func (f *FallbackRateLimitClient) Expire(key string, expiration time.Duration) error {
	f.mutex.RLock()
	isAvailable := f.isRedisAvailable
	f.mutex.RUnlock()

	if isAvailable {
		err := f.redisClient.Expire(key, expiration)
		if err != nil {
			f.switchToFallback()
			return f.memoryStore.Expire(key, expiration)
		}
		return nil
	}

	f.logFallbackWarning()
	return f.memoryStore.Expire(key, expiration)
}

// Delete removes a key
func (f *FallbackRateLimitClient) Delete(key string) error {
	f.mutex.RLock()
	isAvailable := f.isRedisAvailable
	f.mutex.RUnlock()

	if isAvailable {
		err := f.redisClient.Delete(key)
		if err != nil {
			f.switchToFallback()
			return f.memoryStore.Delete(key)
		}
		return nil
	}

	f.logFallbackWarning()
	return f.memoryStore.Delete(key)
}

// switchToFallback marks Redis as unavailable and logs a warning
func (f *FallbackRateLimitClient) switchToFallback() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.isRedisAvailable {
		f.isRedisAvailable = false
		log.Warn().Msg("Redis unavailable, switching to in-memory fallback for rate limiting")
	}
}

// RestoreRedis marks Redis as available again
func (f *FallbackRateLimitClient) RestoreRedis() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if !f.isRedisAvailable {
		f.isRedisAvailable = true
		log.Info().Msg("Redis connection restored, switching back from in-memory fallback")
	}
}

// IsRedisAvailable returns the current Redis availability status
func (f *FallbackRateLimitClient) IsRedisAvailable() bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.isRedisAvailable
}

// logFallbackWarning logs a warning periodically when using fallback
func (f *FallbackRateLimitClient) logFallbackWarning() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	now := time.Now()
	if now.Sub(f.lastWarningTime) >= f.warningInterval {
		log.Warn().Msg("Rate limiting using in-memory fallback (Redis still unavailable)")
		f.lastWarningTime = now
	}
}

// FallbackSlidingWindowClient wraps sliding window operations with fallback
type FallbackSlidingWindowClient struct {
	redisClient      redisClient.SlidingWindowClient
	memoryStore      *InMemorySlidingWindowStore
	isRedisAvailable bool
	mutex            sync.RWMutex
	lastWarningTime  time.Time
	warningInterval  time.Duration
}

// NewFallbackSlidingWindowClient creates a new fallback sliding window client
func NewFallbackSlidingWindowClient(redisClient redisClient.SlidingWindowClient, memoryStore *InMemorySlidingWindowStore) *FallbackSlidingWindowClient {
	return &FallbackSlidingWindowClient{
		redisClient:      redisClient,
		memoryStore:      memoryStore,
		isRedisAvailable: true,
		warningInterval:  5 * time.Minute,
		lastWarningTime:  time.Now(),
	}
}

// ZRemRangeByScore removes entries with scores between min and max
func (f *FallbackSlidingWindowClient) ZRemRangeByScore(ctx context.Context, key, min, max string) error {
	f.mutex.RLock()
	isAvailable := f.isRedisAvailable
	f.mutex.RUnlock()

	if isAvailable {
		err := f.redisClient.ZRemRangeByScore(ctx, key, min, max)
		if err != nil {
			f.switchToFallback()
			return f.memoryStore.ZRemRangeByScore(ctx, key, min, max)
		}
		return nil
	}

	f.logFallbackWarning()
	return f.memoryStore.ZRemRangeByScore(ctx, key, min, max)
}

// ZCount counts entries with scores between min and max
func (f *FallbackSlidingWindowClient) ZCount(ctx context.Context, key, min, max string) (int64, error) {
	f.mutex.RLock()
	isAvailable := f.isRedisAvailable
	f.mutex.RUnlock()

	if isAvailable {
		count, err := f.redisClient.ZCount(ctx, key, min, max)
		if err != nil {
			f.switchToFallback()
			return f.memoryStore.ZCount(ctx, key, min, max)
		}
		return count, nil
	}

	f.logFallbackWarning()
	return f.memoryStore.ZCount(ctx, key, min, max)
}

// ZAdd adds an entry with the given score
func (f *FallbackSlidingWindowClient) ZAdd(ctx context.Context, key string, timestamp int64) error {
	f.mutex.RLock()
	isAvailable := f.isRedisAvailable
	f.mutex.RUnlock()

	if isAvailable {
		err := f.redisClient.ZAdd(ctx, key, timestamp)
		if err != nil {
			f.switchToFallback()
			return f.memoryStore.ZAdd(ctx, key, timestamp)
		}
		return nil
	}

	f.logFallbackWarning()
	return f.memoryStore.ZAdd(ctx, key, timestamp)
}

// Expire sets expiration time on a key
func (f *FallbackSlidingWindowClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	f.mutex.RLock()
	isAvailable := f.isRedisAvailable
	f.mutex.RUnlock()

	if isAvailable {
		err := f.redisClient.Expire(ctx, key, expiration)
		if err != nil {
			f.switchToFallback()
			return f.memoryStore.Expire(ctx, key, expiration)
		}
		return nil
	}

	f.logFallbackWarning()
	return f.memoryStore.Expire(ctx, key, expiration)
}

// switchToFallback marks Redis as unavailable
func (f *FallbackSlidingWindowClient) switchToFallback() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.isRedisAvailable {
		f.isRedisAvailable = false
		log.Warn().Msg("Redis unavailable, switching to in-memory fallback for sliding window rate limiting")
	}
}

// RestoreRedis marks Redis as available again
func (f *FallbackSlidingWindowClient) RestoreRedis() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if !f.isRedisAvailable {
		f.isRedisAvailable = true
		log.Info().Msg("Redis connection restored for sliding window, switching back from in-memory fallback")
	}
}

// IsRedisAvailable returns the current Redis availability status
func (f *FallbackSlidingWindowClient) IsRedisAvailable() bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.isRedisAvailable
}

// logFallbackWarning logs a warning periodically
func (f *FallbackSlidingWindowClient) logFallbackWarning() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	now := time.Now()
	if now.Sub(f.lastWarningTime) >= f.warningInterval {
		log.Warn().Msg("Sliding window rate limiting using in-memory fallback (Redis still unavailable)")
		f.lastWarningTime = now
	}
}
