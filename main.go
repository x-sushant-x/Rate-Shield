package main

import (
	"github.com/x-sushant-x/ThrottleWatch/config"
	"github.com/x-sushant-x/ThrottleWatch/redis"
)

func main() {
	config.LoadConfig()
	redis.Connect()
}
