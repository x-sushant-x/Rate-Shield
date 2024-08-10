package main

import (
	"github.com/x-sushant-x/RateShield/config"
	"github.com/x-sushant-x/RateShield/endpoints"
	"github.com/x-sushant-x/RateShield/limiter"
	"github.com/x-sushant-x/RateShield/redis"
)

func main() {
	config.LoadConfig()
	redis.Connect()

	endpoints.StartTestingRouter()

	limiter.StartSvc()

	select {}
}
