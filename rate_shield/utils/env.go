package utils

import "os"

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

func GetRedisClusterPort() string {
	port := os.Getenv("REDIS_CLUSTER_PORT")

	if len(port) != 0 {
		return port
	}

	return "7001"
}
