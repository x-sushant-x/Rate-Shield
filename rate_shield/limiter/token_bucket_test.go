package limiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/x-sushant-x/RateShield/models"
)

type MockTokenBucketRedis struct {
	mock.Mock
}

func (m *MockTokenBucketRedis) JSONGet(key string) (*models.Bucket, bool, error) {
	args := m.Called(key)

	if bucket, ok := args.Get(0).(*models.Bucket); ok {
		return bucket, args.Bool(1), args.Error(2)
	}

	return nil, args.Bool(1), args.Error(2)
}

func (m *MockTokenBucketRedis) JSONSet(key string, val interface{}) error {
	args := m.Called(key, val)
	return args.Error(0)
}

func (m *MockTokenBucketRedis) Expire(key string, expiration time.Duration) error {
	args := m.Called(key)
	return args.Error(0)
}

func TestTokenBucketService(t *testing.T) {
	mockRedisTokenBucket := new(MockTokenBucketRedis)

	mockRedisTokenBucket.On("JSONSet", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockRedisTokenBucket.On("JSONGet", mock.Anything).Return(nil)
	mockRedisTokenBucket.On("Expire", mock.Anything, mock.Anything).Return(nil)

	tokenBucketSVC := NewTokenBucketService(mockRedisTokenBucket)

	bucket := models.Bucket{
		ClientIP: "192.168.0.1",
		Endpoint: "/api/v1/get-data",
	}
	err := tokenBucketSVC.saveBucket(&bucket)
	assert.NoError(t, err)
}
