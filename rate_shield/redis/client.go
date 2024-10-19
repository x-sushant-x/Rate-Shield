package redisClient

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	TokenBucketClient *redis.Client
	ctx               = context.Background()
)

func createNewRedisConnection(addr string, db int) (*redis.Client, error) {
	conn := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       db,
	})

	pingCmd := conn.Ping(ctx)
	if pingCmd.Err() != nil {
		return nil, pingCmd.Err()
	}
	return conn, nil
}

func NewTokenBucketClient() (RedisTokenBucket, error) {
	client, err := createNewRedisConnection("redis:6379", 1)
	if err != nil {
		return RedisTokenBucket{}, err
	}

	TokenBucketClient = client

	return RedisTokenBucket{
		client: client,
	}, nil
}

func NewFixedWindowClient() (RedisFixedWindow, error) {
	client, err := createNewRedisConnection("redis:6379", 2)
	if err != nil {
		return RedisFixedWindow{}, err
	}

	return RedisFixedWindow{
		client: client,
	}, nil
}

func NewRulesClient() (RedisRules, error) {
	client, err := createNewRedisConnection("redis:6379", 0)
	if err != nil {
		return RedisRules{}, err
	}

	return RedisRules{
		client: client,
	}, nil
}
