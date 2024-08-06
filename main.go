package main

import (
	"github.com/x-sushant-x/RateShield/config"
	"github.com/x-sushant-x/RateShield/redis"
)

func main() {
	config.LoadConfig()
	redis.Connect()
}
