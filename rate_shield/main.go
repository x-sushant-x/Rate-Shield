package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/config"
	"github.com/x-sushant-x/RateShield/endpoints"
	"github.com/x-sushant-x/RateShield/limiter"
	"github.com/x-sushant-x/RateShield/redis"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func main() {
	config.LoadConfig()

	err := redis.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	endpoints.StartTestingRouter()

	limiter.StartSvc()

	select {}
}
