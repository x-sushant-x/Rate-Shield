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

func NewRedisRateLimitClient() (RedisRateLimiterClient, *redis.ClusterClient, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: utils.GetRedisClusterURLs(),
	})

	result, err := client.Ping(ctx).Result()

	if err != nil || result == "" {
		log.Fatal().Err(err).Msg("unable to connect to redis or ping result is nil")
	}

	TokenBucketClient = client

	return RedisRateLimit{
		client: client,
	}, client, nil
}

func NewRulesClient() (RedisRuleClient, error) {
	env := getApplicationEnv()
	port := getRedisRulesInstancePort()

	connectionString := env + port

	client, err := createNewRedisConnection(connectionString, 0)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to redis rules instance on port: " + port)
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
