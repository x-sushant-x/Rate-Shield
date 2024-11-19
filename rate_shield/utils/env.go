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

// (string, string, string) -> URL, Username, Password
func GetRedisRulesInstanceDetails() (string, string) {
	url := os.Getenv("REDIS_RULES_INSTANCE_URL")
	checkEmptyENV(url, "REDIS_RULES_INSTANCE_URL must be provided in docker run command")

	password := os.Getenv("REDIS_RULES_INSTANCE_PASSWORD")

	return url, password
}

func GetRedisClusterURLs() []string {
	clusterURLs := os.Getenv("REDIS_CLUSTERS_URLS")
	checkEmptyENV(clusterURLs, "REDIS_CLUSTERS_URLS not specified in enviroment variables")

	clusterURLsArray := strings.Split(clusterURLs, ",")
	checkEmptyENV(clusterURLs, "REDIS_CLUSTERS_URLS is empty in enviroment variables. Specify comma seperated urls.")

	return clusterURLsArray
}
func GetRedisClusterPassword() string {
	password := os.Getenv("REDIS_CLUSTER_PASSWORD")
	return password
}

func checkEmptyENV(Var string, message string) {
	if len(Var) == 0 {
		log.Fatal().Msg(message)
	}
}
