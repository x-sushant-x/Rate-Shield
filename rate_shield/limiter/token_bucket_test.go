package limiter

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/x-sushant-x/RateShield/models"
)

var (
	redisBucket = new(MockTokenBucketRedis)
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

func TestCreateBucket_Success(t *testing.T) {
	bucket := &models.Bucket{
		ClientIP:        "192.168.0.1",
		Endpoint:        "/api/v1/get-data",
		Capacity:        10,
		AvailableTokens: 10,
		TokenAddRate:    2,
		TokenAddTime:    60,
		CreatedAt:       time.Now().Unix(),
	}

	redisBucket.On("JSONSet", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	redisBucket.On("Expire", mock.Anything, mock.Anything).Return(nil)

	svc := NewTokenBucketService(redisBucket)

	createdBucket, _ := svc.createBucket("192.168.0.1", "/api/v1/get-data", 10, 2)

	assert.Equal(t, bucket, createdBucket)
	redisBucket.AssertExpectations(t)
}

func TestCreateBucket_Fail(t *testing.T) {
	redisBucket.On("JSONSet", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("redis error"))

	svc := NewTokenBucketService(redisBucket)

	createdBucket, err := svc.createBucket("192.168.0.1", "/api/v1/get-data", 10, 2)

	assert.Nil(t, createdBucket)
	assert.NotNil(t, err)
	redisBucket.AssertExpectations(t)
}

func TestCreateBucket_ZeroCapacity(t *testing.T) {
	redisBucket.On("JSONSet", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	redisBucket.On("Expire", mock.Anything, mock.Anything).Return(nil)

	svc := NewTokenBucketService(redisBucket)

	createdBucket, err := svc.createBucket("192.168.0.1", "/api/v1/get-data", 0, 2)

	assert.Nil(t, createdBucket)
	assert.NotNil(t, err)
}

func TestCreateBucket_NegativeTokenAddRate(t *testing.T) {
	redisBucket.On("JSONSet", mock.Anything, mock.Anything).Return(nil)
	redisBucket.On("Expire", mock.Anything, mock.Anything).Return(nil)

	svc := NewTokenBucketService(redisBucket)

	createdBucket, err := svc.createBucket("192.168.0.1", "/api/v1/get-data", 0, -1)

	assert.Nil(t, createdBucket)
	assert.NotNil(t, err)
}
