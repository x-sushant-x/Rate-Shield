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

func NewRedisRateLimitClient() (RedisRateLimiterClient, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"127.0.0.1:6380",
			"127.0.0.1:6381",
			"127.0.0.1:6382",
			"127.0.0.1:6383",
			"127.0.0.1:6384",
			"127.0.0.1:6385",
		},
	})

	result, err := client.Ping(ctx).Result()

	if err != nil || result == "" {
		log.Info().Msg("Error: " + err.Error() + " & Ping Result: " + result)
		return RedisRateLimit{}, err
	}

	TokenBucketClient = client

	return RedisRateLimit{
		client: client,
	}, nil
}

// func NewFixedWindowClient() (RedisFixedWindow, error) {
// 	client, err := createNewRedisConnection(getRedisConnectionStr(), 2)
// 	if err != nil {
// 		return RedisFixedWindow{}, err
// 	}

// 	return RedisFixedWindow{
// 		client: client,
// 	}, nil
// }

// func NewSlidingWindowClient() (*redis.Client, error) {
// 	client, err := createNewRedisConnection(getRedisConnectionStr(), 3)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return client, nil
// }

func NewRulesClient() (RedisRules, error) {
	env := getApplicationEnv()
	port := getRedisRulesInstancePort()

	connectionString := env + port

	client, err := createNewRedisConnection(connectionString, 0)
	if err != nil {
		return RedisRules{}, err
	}

	return RedisRules{
		client: client,
	}, nil
}

func getApplicationEnv() string {
	env := utils.GetApplicationEnviroment()
	redisConnStr := ""

	if env == "prod" {
		redisConnStr = "redis:"
	} else {
		redisConnStr = "localhost:"
	}

	return redisConnStr
}

func getRedisRulesInstancePort() string {
	return utils.GetRedisRulesInstancePort()
}
