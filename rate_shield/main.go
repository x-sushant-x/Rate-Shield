package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/api"
	"github.com/x-sushant-x/RateShield/config"
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
	// Check for configuration file first
	var rulesService service.RulesService
	
	if configPath, exists := config.CheckConfigFileExists("."); exists {
		log.Info().Msgf("Found configuration file: %s", configPath)
		log.Info().Msg("Using file-based rules configuration")
		
		configRulesService, err := service.NewConfigRulesService(configPath)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize config-based rules service")
		}
		rulesService = configRulesService
	} else {
		log.Info().Msg("No configuration file found, using Redis-based rules")
		
		redisRulesClient, err := redisClient.NewRulesClient()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize Redis rules client")
		}
		rulesService = service.NewRedisRulesService(redisRulesClient)
	}

	slackSvc := service.NewSlackService(slackToken, slackChannelID)

	errorNotificationSvc := service.NewErrorNotificationSVC(*slackSvc)

	redisRateLimiter, clusterClient, err := redisClient.NewRedisRateLimitClient()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Redis rate limiter")
	}

	tokenBucketSvc := limiter.NewTokenBucketService(redisRateLimiter, errorNotificationSvc)

	fixedWindowSvc := limiter.NewFixedWindowService(redisRateLimiter)

	slidingWindowSvc := limiter.NewSlidingWindowService(clusterClient)

	limiter := limiter.NewRateLimiterService(&tokenBucketSvc, &fixedWindowSvc, &slidingWindowSvc, rulesService)
	limiter.StartRateLimiter()

	server := api.NewServer(&limiter)
	log.Fatal().Err(server.StartServer())
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
