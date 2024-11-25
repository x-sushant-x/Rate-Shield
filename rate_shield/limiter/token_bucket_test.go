package limiter

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/x-sushant-x/RateShield/models"
	"github.com/x-sushant-x/RateShield/service"
)

type MockRedisRateLimiterClient struct {
	mock.Mock
}

func (m *MockRedisRateLimiterClient) JSONGet(key string) (string, bool, error) {
	args := m.Called(key)
	return args.String(0), args.Bool(1), args.Error(2)
}

func (m *MockRedisRateLimiterClient) JSONSet(key string, val interface{}) error {
	args := m.Called(key, val)
	return args.Error(0)
}

func (m *MockRedisRateLimiterClient) Expire(key string, expireTime time.Duration) error {
	args := m.Called(key, expireTime)
	return args.Error(0)
}

func TestTokenBucketService(t *testing.T) {
	mockRedis := new(MockRedisRateLimiterClient)

	slackSVC := service.NewSlackService("", "")
	errorNotificationSVC := service.NewErrorNotificationSVC(*slackSVC)

	svc := NewTokenBucketService(mockRedis, errorNotificationSVC)

	t.Run("getBucket_success", func(t *testing.T) {
		bucketData := `{"available_tokens" : 10}`

		expectedBucket := &models.Bucket{
			AvailableTokens: 10,
		}

		mockRedis.On("JSONGet", "token_bucket_test").Return(bucketData, true, nil)

		result, found, err := svc.getBucket("test")
		assert.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, expectedBucket, result)

		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil
	})

	t.Run("getBucket_error", func(t *testing.T) {

		mockRedis.On("JSONGet", "token_bucket_test_error").Return("", false, errors.New("redis error"))

		result, found, err := svc.getBucket("test_error")
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.False(t, found)

		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil
	})

	t.Run("getBucket_not_found", func(t *testing.T) {

		mockRedis.On("JSONGet", "token_bucket_test_not_found").Return("", false, nil)

		result, found, err := svc.getBucket("test_not_found")
		assert.Nil(t, result)
		assert.NoError(t, err)
		assert.False(t, found)

		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil
	})

	t.Run("getBucket_unmarshal_error", func(t *testing.T) {
		bucketData := `{"available_tokens" : "10"}`

		mockRedis.On("JSONGet", "token_bucket_test_unmarshal_error").Return(bucketData, true, nil)

		result, found, err := svc.getBucket("test_unmarshal_error")
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.False(t, found)

		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil
	})

	t.Run("saveBucket_success", func(t *testing.T) {
		bucket := &models.Bucket{
			Endpoint:        "/api/v1/get-data",
			Capacity:        100,
			TokenAddRate:    100,
			ClientIP:        "192.168.1.23",
			CreatedAt:       time.Now().Unix(),
			AvailableTokens: 100,
		}

		mockRedis.On("JSONSet", "token_bucket_192.168.1.23:/api/v1/get-data", bucket).Return(nil)
		mockRedis.On("Expire", "token_bucket_192.168.1.23:/api/v1/get-data", time.Second*60).Return(nil)

		err := svc.saveBucket(bucket)
		assert.NoError(t, err)

		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil
	})

	t.Run("saveBucket_success", func(t *testing.T) {
		bucket := &models.Bucket{
			Endpoint:        "/api/v1/get-data",
			Capacity:        100,
			TokenAddRate:    100,
			ClientIP:        "192.168.1.23",
			CreatedAt:       time.Now().Unix(),
			AvailableTokens: 100,
		}

		mockRedis.On("JSONSet", "token_bucket_192.168.1.23:/api/v1/get-data", bucket).Return(nil)
		mockRedis.On("Expire", "token_bucket_192.168.1.23:/api/v1/get-data", time.Second*60).Return(nil)

		err := svc.saveBucket(bucket)
		assert.NoError(t, err)

		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil
	})

	t.Run("saveBucket_error", func(t *testing.T) {
		bucket := &models.Bucket{
			Endpoint:        "/api/v1/get-data",
			Capacity:        100,
			TokenAddRate:    100,
			ClientIP:        "192.168.1.23",
			CreatedAt:       time.Now().Unix(),
			AvailableTokens: 100,
		}

		mockRedis.On("JSONSet", "token_bucket_192.168.1.23:/api/v1/get-data", bucket).Return(errors.New("redis-error"))
		mockRedis.On("Expire", "token_bucket_192.168.1.23:/api/v1/get-data", time.Second*60).Return(nil)

		err := svc.saveBucket(bucket)
		assert.Error(t, err)

		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil
	})

	t.Run("saveBucket_expire_error", func(t *testing.T) {
		bucket := &models.Bucket{
			Endpoint:        "/api/v1/get-data",
			Capacity:        100,
			TokenAddRate:    100,
			ClientIP:        "192.168.1.23",
			CreatedAt:       time.Now().Unix(),
			AvailableTokens: 100,
		}

		mockRedis.On("JSONSet", "token_bucket_192.168.1.23:/api/v1/get-data", bucket).Return(nil)
		mockRedis.On("Expire", "token_bucket_192.168.1.23:/api/v1/get-data", time.Second*60).Return(errors.New("redis-error"))

		err := svc.saveBucket(bucket)
		assert.Error(t, err)

		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil
	})

	t.Run("parseKey_success", func(t *testing.T) {
		_, _, err := parseKey("192.168.23.1:/api/v1/get-data")
		assert.NoError(t, err)
	})

	t.Run("parseKey_error", func(t *testing.T) {
		_, _, err := parseKey("")
		assert.Error(t, err)
	})
}
