package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/api"
	"github.com/x-sushant-x/RateShield/limiter"
	redisClient "github.com/x-sushant-x/RateShield/redis"
	"github.com/x-sushant-x/RateShield/service"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func main() {
	redisTokenBucket, err := redisClient.NewTokenBucketClient()
	if err != nil {
		log.Fatal().Err(err)
	}

	redisFixedWindow, err := redisClient.NewFixedWindowClient()
	if err != nil {
		log.Fatal().Err(err)
	}

	redisRulesClient, err := redisClient.NewRulesClient()
	if err != nil {
		log.Fatal().Err(err)
	}

	// TODO -> Utilize this service later.
	// TODO -> Add slackToken and channel in env.
	_ = service.NewSlackService("slackToken", "channel")

	tokenBucketSvc := limiter.NewTokenBucketService(redisTokenBucket)
	fixedWindowSvc := limiter.NewFixedWindowService(redisFixedWindow)
	redisRulesSvc := service.NewRedisRulesService(redisRulesClient)

	limiter := limiter.NewRateLimiterService(&tokenBucketSvc, &fixedWindowSvc, redisRulesSvc)
	limiter.StartRateLimiter()

	server := api.NewServer(8080)
	log.Fatal().Err(server.StartServer())

	select {}
}
