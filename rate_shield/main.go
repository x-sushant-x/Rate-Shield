package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/api"
	"github.com/x-sushant-x/RateShield/fallback"
	"github.com/x-sushant-x/RateShield/limiter"
	redisClient "github.com/x-sushant-x/RateShield/redis"
	"github.com/x-sushant-x/RateShield/service"
	"github.com/x-sushant-x/RateShield/utils"
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
	fallbackEnabled := utils.GetRedisFallbackEnabled()

	if fallbackEnabled {
		log.Info().Msg("Redis fallback is ENABLED - will use in-memory storage if Redis is unavailable")
	} else {
		log.Info().Msg("Redis fallback is DISABLED - application will fail if Redis is unavailable")
	}

	redisRulesClient, err := redisClient.NewRulesClient()
	if err != nil {
		if !fallbackEnabled {
			log.Fatal().Err(err)
		}
		log.Warn().Msg("Redis rules client unavailable, continuing with fallback")
	}

	slackSvc := service.NewSlackService(slackToken, slackChannelID)

	errorNotificationSvc := service.NewErrorNotificationSVC(*slackSvc)

	redisRateLimiter, clusterClient, err := redisClient.NewRedisRateLimitClient()
	if err != nil {
		if !fallbackEnabled {
			log.Fatal().Err(err)
		}
		log.Warn().Msg("Redis rate limit client unavailable, initializing in-memory fallback")

		// Initialize in-memory stores
		memoryStore := fallback.NewInMemoryRateLimitStore()
		memorySWStore := fallback.NewInMemorySlidingWindowStore()

		// Create fallback clients with nil Redis clients
		rateLimitClient := fallback.NewFallbackRateLimitClient(nil, memoryStore)
		slidingWindowClient := fallback.NewFallbackSlidingWindowClient(nil, memorySWStore)

		// Mark as unavailable from the start
		rateLimitClient.RestoreRedis() // This won't actually restore, just sets the state

		// Create services with fallback clients
		tokenBucketSvc := limiter.NewTokenBucketService(rateLimitClient, errorNotificationSvc)
		fixedWindowSvc := limiter.NewFixedWindowService(rateLimitClient)

		redisRulesSvc := service.NewRedisRulesService(redisRulesClient)

		slidingWindowSvc := limiter.NewSlidingWindowService(slidingWindowClient)

		limiter := limiter.NewRateLimiterService(&tokenBucketSvc, &fixedWindowSvc, &slidingWindowSvc, redisRulesSvc)
		limiter.StartRateLimiter()

		// Start health monitor to check for Redis recovery
		if clusterClient != nil {
			retryInterval := utils.GetRedisRetryInterval()
			healthMonitor := fallback.NewRedisHealthMonitor(clusterClient, rateLimitClient, slidingWindowClient, retryInterval)
			healthMonitor.Start()
		}

		server := api.NewServer(&limiter)
		log.Fatal().Err(server.StartServer())
		return
	}

	// Normal flow when Redis is available
	tokenBucketSvc := limiter.NewTokenBucketService(redisRateLimiter, errorNotificationSvc)
	fixedWindowSvc := limiter.NewFixedWindowService(redisRateLimiter)
	redisRulesSvc := service.NewRedisRulesService(redisRulesClient)

	// Create sliding window client with interface
	slidingWindowClient := redisClient.NewSlidingWindowClient(clusterClient)

	// If fallback is enabled, wrap with fallback client
	var finalSWClient redisClient.SlidingWindowClient = slidingWindowClient
	var finalRateLimitClient redisClient.RedisRateLimiterClient = redisRateLimiter

	if fallbackEnabled {
		log.Info().Msg("Wrapping Redis clients with fallback support")
		memoryStore := fallback.NewInMemoryRateLimitStore()
		memorySWStore := fallback.NewInMemorySlidingWindowStore()

		rateLimitFallback := fallback.NewFallbackRateLimitClient(redisRateLimiter, memoryStore)
		slidingWindowFallback := fallback.NewFallbackSlidingWindowClient(slidingWindowClient, memorySWStore)

		finalRateLimitClient = rateLimitFallback
		finalSWClient = slidingWindowFallback

		// Start health monitor
		retryInterval := utils.GetRedisRetryInterval()
		healthMonitor := fallback.NewRedisHealthMonitor(clusterClient, rateLimitFallback, slidingWindowFallback, retryInterval)
		healthMonitor.Start()

		// Recreate services with fallback clients
		tokenBucketSvc = limiter.NewTokenBucketService(finalRateLimitClient, errorNotificationSvc)
		fixedWindowSvc = limiter.NewFixedWindowService(finalRateLimitClient)
	}

	slidingWindowSvc := limiter.NewSlidingWindowService(finalSWClient)

	limiter := limiter.NewRateLimiterService(&tokenBucketSvc, &fixedWindowSvc, &slidingWindowSvc, redisRulesSvc)
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
