package redisClient

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisSlidingWindowClient wraps redis.ClusterClient to implement SlidingWindowClient interface
type RedisSlidingWindowClient struct {
	client *redis.ClusterClient
}

// NewRedisSlidingWindowClient creates a new Redis sliding window client
func NewRedisSlidingWindowClient(client *redis.ClusterClient) *RedisSlidingWindowClient {
	return &RedisSlidingWindowClient{
		client: client,
	}
}

// ZRemRangeByScore removes entries with scores between min and max
func (r *RedisSlidingWindowClient) ZRemRangeByScore(ctx context.Context, key, min, max string) error {
	return r.client.ZRemRangeByScore(ctx, key, min, max).Err()
}

// ZCount counts entries with scores between min and max
func (r *RedisSlidingWindowClient) ZCount(ctx context.Context, key, min, max string) (int64, error) {
	return r.client.ZCount(ctx, key, min, max).Result()
}

// ZAdd adds an entry with the given score
func (r *RedisSlidingWindowClient) ZAdd(ctx context.Context, key string, timestamp int64) error {
	return r.client.ZAdd(ctx, key, redis.Z{
		Member: timestamp,
		Score:  float64(timestamp),
	}).Err()
}

// Expire sets expiration time on a key
func (r *RedisSlidingWindowClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// ExecutePipeline executes a pipeline of commands for sliding window operations
func (r *RedisSlidingWindowClient) ExecutePipeline(ctx context.Context, key string, now int64, windowSize time.Duration) (int64, error) {
	then := fmt.Sprintf("%d", now-int64(windowSize.Seconds()))

	pipe := r.client.TxPipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", then)
	countCmd := pipe.ZCount(ctx, key, then, fmt.Sprintf("%d", now))

	_, err := pipe.Exec(ctx)
	if err != nil {
		return -1, err
	}

	count, err := countCmd.Result()
	if err != nil {
		return -1, err
	}

	return count, nil
}

// UpdateWindow adds a timestamp and sets expiration in a pipeline
func (r *RedisSlidingWindowClient) UpdateWindow(ctx context.Context, key string, now int64, windowSize time.Duration) error {
	pipe := r.client.TxPipeline()

	pipe.ZAdd(ctx, key, redis.Z{
		Member: now,
		Score:  float64(now),
	})

	pipe.Expire(ctx, key, windowSize)

	_, err := pipe.Exec(ctx)
	return err
}
