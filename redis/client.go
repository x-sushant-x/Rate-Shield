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

func SetJSONObject(key string, val interface{}) error {
	ctx := context.Background()
	err := Client.JSONSet(ctx, key, ".", val).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetJSONObject(key string) ([]byte, error) {
	ctx := context.Background()
	res, err := Client.JSONGet(ctx, key, ".").Result()
	if err != nil {
		return nil, err
	}
	return []byte(res), nil
}
