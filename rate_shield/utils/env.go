package utils

import (
	"os"
	"strconv"
	"strings"
	"time"

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

// GetRedisFallbackEnabled returns whether in-memory fallback is enabled when Redis is unavailable
func GetRedisFallbackEnabled() bool {
	fallbackEnabled := os.Getenv("ENABLE_REDIS_FALLBACK")
	if len(fallbackEnabled) == 0 {
		return false
	}

	enabled, err := strconv.ParseBool(fallbackEnabled)
	if err != nil {
		log.Warn().Msgf("Invalid ENABLE_REDIS_FALLBACK value: %s, defaulting to false", fallbackEnabled)
		return false
	}

	return enabled
}

// GetRedisRetryInterval returns the interval for retrying Redis connection when using fallback
func GetRedisRetryInterval() time.Duration {
	retryInterval := os.Getenv("REDIS_RETRY_INTERVAL")
	if len(retryInterval) == 0 {
		return 30 * time.Second // default 30 seconds
	}

	duration, err := time.ParseDuration(retryInterval)
	if err != nil {
		log.Warn().Msgf("Invalid REDIS_RETRY_INTERVAL value: %s, defaulting to 30s", retryInterval)
		return 30 * time.Second
	}

	return duration
}

func checkEmptyENV(Var string, message string) {
	if len(Var) == 0 {
		log.Fatal().Msg(message)
	}
}
