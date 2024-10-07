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
	"log"
	"rate_shield/service"
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

	// Get the Slack token from environment variable
	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		log.Fatal("SLACK_TOKEN environment variable is not set")
	}

	// The channel you want to send messages to
	channel := "#alerts" // Replace with your Slack channel name or channel ID

	// Create a new Slack service
	slackService := service.NewSlackService(slackToken, channel)

	// Send a message to Slack
	message := "Rate Shield has failed, please investigate."
	err := slackService.SendSlackMessage(message)
	if err != nil {
		log.Fatalf("Error sending Slack message: %v", err)
	}

	log.Println("Slack message sent successfully")

	tokenBucketSvc := limiter.NewTokenBucketService(redisTokenBucket)
	fixedWindowSvc := limiter.NewFixedWindowService(redisFixedWindow)
	redisRulesSvc := service.NewRedisRulesService(redisRulesClient)

	limiter := limiter.NewRateLimiterService(&tokenBucketSvc, &fixedWindowSvc, redisRulesSvc)
	limiter.StartRateLimiter()

	server := api.NewServer(8080)
	log.Fatal().Err(server.StartServer())

	select {}
}
