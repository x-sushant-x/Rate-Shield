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
		log.Info().Msg("Redis fallback is ENABLED for rate limit cluster only (rules instance always required)")
	} else {
		log.Info().Msg("Redis fallback is DISABLED")
	}

	// ========================================
	// Rules Redis Instance - ALWAYS REQUIRED
	// ========================================
	log.Info().Msg("Connecting to Rules Redis instance (required)...")
	redisRulesClient, rulesInstanceClient, err := redisClient.NewRulesClient()
	if err != nil {
		log.Fatal().Err(err).Msg("Rules Redis instance is required - cannot start without it")
	}
	log.Info().Msg("Rules Redis instance connected successfully ✅")

	slackSvc := service.NewSlackService(slackToken, slackChannelID)
	errorNotificationSvc := service.NewErrorNotificationSVC(*slackSvc)
	redisRulesSvc := service.NewRedisRulesService(redisRulesClient)

	// ========================================
	// Rate Limit Cluster - Supports Fallback
	// ========================================
	log.Info().Msg("Connecting to Rate Limit Redis Cluster...")
	redisRateLimiter, clusterClient, err := redisClient.NewRedisRateLimitClient()
	if err != nil {
		if !fallbackEnabled {
			log.Fatal().Err(err).Msg("Rate Limit Redis Cluster is required (fallback disabled)")
		}
		log.Warn().Msg("Rate Limit Redis Cluster unavailable, initializing in-memory fallback")

		// Initialize in-memory stores
		memoryStore := fallback.NewInMemoryRateLimitStore()
		memorySWStore := fallback.NewInMemorySlidingWindowStore()

		// Create fallback clients with nil Redis clients
		rateLimitClient := fallback.NewFallbackRateLimitClient(nil, memoryStore)
		slidingWindowClient := fallback.NewFallbackSlidingWindowClient(nil, memorySWStore)

		// Create services with fallback clients
		tokenBucketSvc := limiter.NewTokenBucketService(rateLimitClient, errorNotificationSvc)
		fixedWindowSvc := limiter.NewFixedWindowService(rateLimitClient)
		slidingWindowSvc := limiter.NewSlidingWindowService(slidingWindowClient)

		rateLimiter := limiter.NewRateLimiterService(&tokenBucketSvc, &fixedWindowSvc, &slidingWindowSvc, redisRulesSvc)
		rateLimiter.StartRateLimiter()

		// Start rules health monitor
		retryInterval := utils.GetRedisRetryInterval()
		rulesHealthMonitor := fallback.NewRulesHealthMonitor(rulesInstanceClient, &rateLimiter, retryInterval)
		rulesHealthMonitor.Start()

		log.Warn().Msg("Rate Limit Cluster health monitoring unavailable (cluster failed at startup)")

		server := api.NewServer(&rateLimiter)
		log.Fatal().Err(server.StartServer())
		return
	}
	log.Info().Msg("Rate Limit Redis Cluster connected successfully ✅")

	// ========================================
	// Both Redis Instances Available
	// ========================================
	slidingWindowClient := redisClient.NewSlidingWindowClient(clusterClient)

	// Determine if we should wrap with fallback support
	var finalSWClient redisClient.SlidingWindowClient = slidingWindowClient
	var finalRateLimitClient redisClient.RedisRateLimiterClient = redisRateLimiter

	if fallbackEnabled {
		log.Info().Msg("Wrapping Rate Limit Cluster clients with fallback support")
		memoryStore := fallback.NewInMemoryRateLimitStore()
		memorySWStore := fallback.NewInMemorySlidingWindowStore()

		rateLimitFallback := fallback.NewFallbackRateLimitClient(redisRateLimiter, memoryStore)
		slidingWindowFallback := fallback.NewFallbackSlidingWindowClient(slidingWindowClient, memorySWStore)

		finalRateLimitClient = rateLimitFallback
		finalSWClient = slidingWindowFallback

		// Start health monitors
		retryInterval := utils.GetRedisRetryInterval()

		// Monitor rate limit cluster
		rateLimitHealthMonitor := fallback.NewRedisHealthMonitor(clusterClient, rateLimitFallback, slidingWindowFallback, retryInterval)
		rateLimitHealthMonitor.Start()

		// Monitor rules instance
		rulesHealthMonitor := fallback.NewRulesHealthMonitor(rulesInstanceClient, nil, retryInterval)
		rulesHealthMonitor.Start()
	}

	// Create services
	tokenBucketSvc := limiter.NewTokenBucketService(finalRateLimitClient, errorNotificationSvc)
	fixedWindowSvc := limiter.NewFixedWindowService(finalRateLimitClient)
	slidingWindowSvc := limiter.NewSlidingWindowService(finalSWClient)

	rateLimiter := limiter.NewRateLimiterService(&tokenBucketSvc, &fixedWindowSvc, &slidingWindowSvc, redisRulesSvc)
	rateLimiter.StartRateLimiter()

	// Start rules health monitor (always, regardless of fallback setting)
	if fallbackEnabled {
		retryInterval := utils.GetRedisRetryInterval()
		rulesHealthMonitor := fallback.NewRulesHealthMonitor(rulesInstanceClient, &rateLimiter, retryInterval)
		rulesHealthMonitor.Start()
	}

	server := api.NewServer(&rateLimiter)
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
