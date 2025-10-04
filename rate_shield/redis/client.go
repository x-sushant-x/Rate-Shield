package redisClient

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/utils"
)

var (
	ctx               = context.Background()
	TokenBucketClient *redis.ClusterClient
)

func createNewRedisConnection(addr, password string) (*redis.Client, error) {
	conn := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	result, err := conn.Ping(ctx).Result()
	if err != nil || result == "" {
		// Rules instance is ALWAYS required - no fallback
		log.Fatal().Err(err).Msg("unable to connect to redis rules instance: " + addr)
	}

	return conn, nil
}

func NewRedisRateLimitClient() (RedisRateLimiterClient, *redis.ClusterClient, error) {
	clusterURLs := utils.GetRedisClusterURLs()
	clusterPassword := utils.GetRedisClusterPassword()

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Password: clusterPassword,
		Addrs:    clusterURLs,
	})

	result, err := client.Ping(ctx).Result()
	if err != nil || result == "" {
		// Return error instead of fatal if fallback is enabled
		if utils.GetRedisFallbackEnabled() {
			log.Warn().Err(err).Msg("unable to connect to redis rate limit cluster (fallback enabled, will use in-memory storage)")
			return nil, nil, err
		}
		log.Fatal().Err(err).Msg("unable to connect to redis or ping result is nil for rate limit cluster")
	}

	TokenBucketClient = client

	return RedisRateLimit{
		client: client,
	}, client, nil
}

func NewRulesClient() (RedisRuleClient, *redis.Client, error) {
	url, password := utils.GetRedisRulesInstanceDetails()

	client, err := createNewRedisConnection(url, password)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to redis rules instance: " + url)
	}

	return RedisRules{
		client: client,
	}, client, nil
}

func NewSlidingWindowClient(clusterClient *redis.ClusterClient) SlidingWindowClient {
	return NewRedisSlidingWindowClient(clusterClient)
}
