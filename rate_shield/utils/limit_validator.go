package utils

func ValidateLimitRequest(ip, endpoint string) error {
	if len(ip) == 0 {
		return ErrorInvalidIP
	}

	if len(endpoint) == 0 {
		return ErrorInvalidEndpoint
	}

	return nil
}
