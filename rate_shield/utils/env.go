package utils

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

func GetApplicationEnviroment() string {
	appEnv := os.Getenv("ENV")

	if len(appEnv) != 0 {
		return appEnv
	}

	return "dev"
}

func GetRedisRulesInstancePort() string {
	port := os.Getenv("REDIS_RULES_INSTANCE_PORT")

	if len(port) != 0 {
		return port
	}

	return "7000"
}

func GetRedisClusterURLs() []string {
	clusterURLs := os.Getenv("REDIS_CLUSTERS_URLS")
	if len(clusterURLs) == 0 {
		log.Fatal().Msg("REDIS_CLUSTERS_URLS not specified in enviroment variables")
	}

	clusterURLsArray := strings.Split(clusterURLs, ",")
	if len(clusterURLsArray) == 0 {
		log.Fatal().Msg("REDIS_CLUSTERS_URLS is empty in enviroment variables. Specify comma seperated urls.")
	}

	return clusterURLsArray
}
