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

	redisRulesClient, err := redisClient.NewRulesClient()
	if err != nil {
		log.Fatal().Err(err)
	}

	authSvc := service.NewAuthService(redisClient.RuleClientConn)
	authSvc.InitializeDefaultCreds()

	// Create audit client and service
	auditClient := redisClient.NewAuditClient(redisRulesClient.(redisClient.RedisRules).GetClient())
	auditSvc := service.NewAuditService(auditClient)

	slackSvc := service.NewSlackService(slackToken, slackChannelID)

	errorNotificationSvc := service.NewErrorNotificationSVC(*slackSvc)

	redisRateLimiter, clusterClient, err := redisClient.NewRedisRateLimitClient()
	if err != nil {
		log.Fatal().Err(err)
	}

	tokenBucketSvc := limiter.NewTokenBucketService(redisRateLimiter, errorNotificationSvc)

	fixedWindowSvc := limiter.NewFixedWindowService(redisRateLimiter)

	redisRulesSvc := service.NewRedisRulesService(redisRulesClient, auditSvc)

	slidingWindowSvc := limiter.NewSlidingWindowService(clusterClient)

	limiter := limiter.NewRateLimiterService(&tokenBucketSvc, &fixedWindowSvc, &slidingWindowSvc, redisRulesSvc)
	limiter.StartRateLimiter()

	go func() {
		server := api.NewServer(&limiter)
		log.Fatal().Err(server.StartServer())
	}()

	go func() {
		api.StartGRPCServer(&limiter, "50051")
	}()

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
