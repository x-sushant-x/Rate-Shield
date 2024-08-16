package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var (
	Client *redis.Client
	ctx    = context.Background()
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
		log.Error().Msgf("unable to connect to redis: " + cmd.Err().Error())
	} else {
		log.Info().Msg("Connected to redis successfully")
	}
}

func SetJSONObject(key string, val interface{}) error {
	err := Client.JSONSet(ctx, key, ".", val).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetJSONObject(key string) ([]byte, error) {
	res, err := Client.JSONGet(ctx, key, ".").Result()
	if err != nil {
		return nil, err
	}
	return []byte(res), nil
}

func Get(key string) ([]byte, bool, error) {
	res, err := Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	return []byte(res), true, nil
}
