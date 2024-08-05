package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
)

func Connect() {
	c := redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
	)

	Client = c

	cmd := Client.Ping(context.TODO())
	if cmd.Err() != nil {
		log.Fatal("can not connect to redis: " + cmd.Err().Error())
	}
}
