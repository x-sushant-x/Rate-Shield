package config

import (
	"log"

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
		log.Fatal("unable to read config: " + err.Error())
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Fatal("unable to marshal config: " + err.Error())
	}
}
