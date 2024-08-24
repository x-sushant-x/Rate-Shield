package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/api"
	"github.com/x-sushant-x/RateShield/config"
	"github.com/x-sushant-x/RateShield/limiter"
	redisClient "github.com/x-sushant-x/RateShield/redis"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func main() {
	config.LoadConfig()

	err := redisClient.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	tokenBucketSvc := limiter.NewTokenBucketService()
	limiter := limiter.NewRateLimiterService(tokenBucketSvc)
	limiter.StartRateLimiter()

	server := api.NewServer(8080)
	log.Fatal().Err(server.StartServer())

	select {}
}
