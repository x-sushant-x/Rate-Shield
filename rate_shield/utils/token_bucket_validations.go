package utils

func ValidateCreateBucketReq(ip, endpoint string, capacity, tokenAddRate int) error {
	if len(ip) == 0 {
		return ErrorInvalidIP
	}

	if len(endpoint) == 0 {
		return ErrorInvalidEndpoint
	}

	if capacity <= 0 {
		return ErrorZeroCapacity

	}

	if tokenAddRate <= 0 {
		return ErrorNegativeAddTokenRate
	}

	return nil
}
