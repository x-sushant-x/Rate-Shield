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

func TestSpawnNewFixedWindow(t *testing.T) {
	mockRedis := new(MockRedisFixedWindowClient)
	service := NewFixedWindowService(mockRedis)

	t.Run("Successful new window creation", func(t *testing.T) {
		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil

		// Mock the JSONSet operation
		mockRedis.On("JSONSet", mock.Anything, mock.Anything).Return(nil)

		// Mock the Expire operation
		mockRedis.On("Expire", mock.Anything, mock.Anything).Return(nil)

		rule := &models.Rule{
			FixedWindowCounterRule: &models.FixedWindowCounterRule{
				MaxRequests: 10,
				Window:      60,
			},
		}

		fixedWindow, err := service.spawnNewFixedWindow("192.168.1.6", "/test", rule)

		// Assert that no error occurred
		assert.NoError(t, err)

		// Assert that the fixedWindow is not nil
		assert.NotNil(t, fixedWindow)

		// Assert the properties of the created fixedWindow
		assert.Equal(t, "/test", fixedWindow.Endpoint)
		assert.Equal(t, "192.168.1.6", fixedWindow.ClientIP)
		assert.Equal(t, int64(10), fixedWindow.MaxRequests)
		assert.Equal(t, int64(1), fixedWindow.CurrRequests)
		assert.Equal(t, 60, fixedWindow.Window)

		// Assert that the mock expectations were met
		mockRedis.AssertExpectations(t)
	})

	t.Run("Redis JSONSet error", func(t *testing.T) {
		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil

		// Mock the JSONSet operation to return an error
		expectedError := errors.New("Redis JSONSet error")
		mockRedis.On("JSONSet", mock.Anything, mock.Anything).Return(expectedError)

		rule := &models.Rule{
			FixedWindowCounterRule: &models.FixedWindowCounterRule{
				MaxRequests: 10,
				Window:      60,
			},
		}

		fixedWindow, err := service.spawnNewFixedWindow("192.168.1.7", "/test", rule)

		// Assert that an error occurred
		assert.Error(t, err)
		assert.Nil(t, fixedWindow)

		// Check if the error is exactly what we expect
		assert.Equal(t, expectedError, err)

		// Assert that the mock expectations were met
		mockRedis.AssertExpectations(t)
	})

	t.Run("Redis Expire error", func(t *testing.T) {
		mockRedis.ExpectedCalls = nil
		mockRedis.Calls = nil

		// Mock the JSONSet operation to succeed
		mockRedis.On("JSONSet", mock.Anything, mock.Anything).Return(nil)

		// Mock the Expire operation to fail
		expectedError := errors.New("Redis Expire error")
		mockRedis.On("Expire", mock.Anything, mock.Anything).Return(expectedError)

		rule := &models.Rule{
			FixedWindowCounterRule: &models.FixedWindowCounterRule{
				MaxRequests: 10,
				Window:      60,
			},
		}

		fixedWindow, err := service.spawnNewFixedWindow("192.168.1.8", "/test", rule)

		// Assert that an error occurred
		assert.Error(t, err)
		assert.Nil(t, fixedWindow)

		// Check if the error is exactly what we expect
		assert.Equal(t, expectedError, err)

		// Assert that the mock expectations were met
		mockRedis.AssertExpectations(t)
	})
}
