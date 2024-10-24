package utils

import "os"

func GetApplicationEnviroment() string {
	appEnv := os.Getenv("ENV")

	if len(appEnv) != 0 {
		return appEnv
	}

	return "dev"
}
