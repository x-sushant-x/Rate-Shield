package utils

import "errors"

var (
	ErrorNegativeAddTokenRate = errors.New("invalid token add rate. Must be greater than 0")
	ErrorZeroCapacity         = errors.New("invalid token capacity. Must be greater than 0")
)

var (
	ErrorInvalidIP       = errors.New("invalid IP Address. Make sure it's not empty")
	ErrorInvalidEndpoint = errors.New("invalid API Endpoint. Make sure it's not empty")
)
