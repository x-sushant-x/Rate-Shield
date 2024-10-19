package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/api"
	"github.com/x-sushant-x/RateShield/limiter"
	redisClient "github.com/x-sushant-x/RateShield/redis"
	"github.com/x-sushant-x/RateShield/service"
)

var (
	slackToken     string
	slackChannelID string
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	loadENVFile()
	setSlackCredentials()
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

	slackSvc := service.NewSlackService(slackToken, slackChannelID)

	errorNotificationSvc := service.NewErrorNotificationSVC(*slackSvc)

	tokenBucketSvc := limiter.NewTokenBucketService(redisTokenBucket, errorNotificationSvc)

	fixedWindowSvc := limiter.NewFixedWindowService(redisFixedWindow)

	redisRulesSvc := service.NewRedisRulesService(redisRulesClient)

	limiter := limiter.NewRateLimiterService(&tokenBucketSvc, &fixedWindowSvc, redisRulesSvc)
	limiter.StartRateLimiter()

	server := api.NewServer(8080)
	log.Fatal().Err(server.StartServer())

	select {}
}

func loadENVFile() {
	err := godotenv.Load()
	if err != nil {
		log.Panic().Msgf("error while loading env file: %s", err)
	}
}

func setSlackCredentials() {
	sToken := os.Getenv("SLACK_TOKEN")
	if len(sToken) == 0 {
		log.Panic().Msg("SLACK_TOKEN not available in env file")
	}
	slackToken = sToken

	sChannel := os.Getenv("SLACK_CHANNEL")
	if len(sChannel) == 0 {
		log.Panic().Msg("SLACK_CHANNEL not available in env file")
	}
	slackChannelID = sChannel
}
