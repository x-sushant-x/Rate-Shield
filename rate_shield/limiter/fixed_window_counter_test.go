package limiter

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/x-sushant-x/RateShield/models"
)

type MockRedisFixedWindowClient struct {
	mock.Mock
}

func (m *MockRedisFixedWindowClient) JSONGet(key string) (*models.FixedWindowCounter, bool, error) {
	args := m.Called(key)
	return args.Get(0).(*models.FixedWindowCounter), args.Bool(1), args.Error(2)
}

func (m *MockRedisFixedWindowClient) JSONSet(key string, value interface{}) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockRedisFixedWindowClient) Expire(key string, expiration time.Duration) error {
	args := m.Called(key, expiration)
	return args.Error(0)
}

func TestProcessRequest(t *testing.T) {
	mockRedis := new(MockRedisFixedWindowClient)
	service := NewFixedWindowService(mockRedis)

	t.Run("New window creation", func(t *testing.T) {
		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil

		mockRedis.On("JSONGet", mock.Anything).Return((*models.FixedWindowCounter)(nil), false, nil)
		mockRedis.On("JSONSet", mock.Anything, mock.Anything).Return(nil)
		mockRedis.On("Expire", mock.Anything, mock.Anything).Return(nil)

		rule := &models.Rule{
			FixedWindowCounterRule: &models.FixedWindowCounterRule{
				MaxRequests: 10,
				Window:      60,
			},
		}

		response := service.processRequest("192.168.1.1", "/test", rule)

		assert.Equal(t, 200, response.HTTPStatusCode)
		mockRedis.AssertExpectations(t)
	})

	t.Run("Existing window within limit", func(t *testing.T) {
		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil

		mockRedis.On("JSONGet", mock.Anything).Return(&models.FixedWindowCounter{
			MaxRequests:    10,
			CurrRequests:   5,
			Window:         60,
			LastAccessTime: time.Now().Unix() - 30,
		}, true, nil)
		mockRedis.On("JSONSet", mock.Anything, mock.Anything).Return(nil)

		rule := &models.Rule{
			FixedWindowCounterRule: &models.FixedWindowCounterRule{
				MaxRequests: 10,
				Window:      60,
			},
		}

		response := service.processRequest("192.168.1.2", "/test", rule)

		assert.Equal(t, 200, response.HTTPStatusCode)
		mockRedis.AssertExpectations(t)
	})

	t.Run("Existing window at limit", func(t *testing.T) {
		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil

		mockRedis.On("JSONGet", mock.Anything).Return(&models.FixedWindowCounter{
			MaxRequests:    10,
			CurrRequests:   10,
			Window:         60,
			LastAccessTime: time.Now().Unix() - 30,
		}, true, nil)

		rule := &models.Rule{
			FixedWindowCounterRule: &models.FixedWindowCounterRule{
				MaxRequests: 10,
				Window:      60,
			},
		}

		response := service.processRequest("192.168.1.3", "/test", rule)

		assert.Equal(t, 429, response.HTTPStatusCode)
		mockRedis.AssertExpectations(t)
	})

	t.Run("Existing window outside time limit", func(t *testing.T) {
		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil

		mockRedis.On("JSONGet", mock.Anything).Return(&models.FixedWindowCounter{
			MaxRequests:    10,
			CurrRequests:   10,
			Window:         60,
			LastAccessTime: time.Now().Unix() - 61,
		}, true, nil)
		mockRedis.On("JSONSet", mock.Anything, mock.Anything).Return(nil)

		rule := &models.Rule{
			FixedWindowCounterRule: &models.FixedWindowCounterRule{
				MaxRequests: 10,
				Window:      60,
			},
		}

		response := service.processRequest("192.168.1.4", "/test", rule)

		assert.Equal(t, 200, response.HTTPStatusCode)
		mockRedis.AssertExpectations(t)
	})

	t.Run("Redis error", func(t *testing.T) {
		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil

		mockRedis.On("JSONGet", mock.Anything).Return((*models.FixedWindowCounter)(nil), false, errors.New("Redis error"))

		rule := &models.Rule{
			FixedWindowCounterRule: &models.FixedWindowCounterRule{
				MaxRequests: 10,
				Window:      60,
			},
		}

		response := service.processRequest("192.168.1.5", "/test", rule)

		assert.Equal(t, 500, response.HTTPStatusCode)
		mockRedis.AssertExpectations(t)
	})
}
