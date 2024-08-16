package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	Config configuration
)

type configuration struct {
	TokenAddingRate     int `json:"tokenAddingRate"`
	TokenBucketCapacity int `json:"tokenBucketCapacity"`
}

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Error().Msgf("error reading config file: %s", err.Error())
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Error().Msgf("error while marshal config: %s" + err.Error())
	}
}
