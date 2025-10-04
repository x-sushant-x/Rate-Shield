package fallback

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// RedisHealthMonitor monitors Redis connection health and manages reconnection
type RedisHealthMonitor struct {
	rateLimitClient  *redis.ClusterClient
	fallbackClient   *FallbackRateLimitClient
	fallbackSWClient *FallbackSlidingWindowClient
	retryInterval    time.Duration
	stopCh           chan struct{}
}

// NewRedisHealthMonitor creates a new health monitor
func NewRedisHealthMonitor(
	rateLimitClient *redis.ClusterClient,
	fallbackClient *FallbackRateLimitClient,
	fallbackSWClient *FallbackSlidingWindowClient,
	retryInterval time.Duration,
) *RedisHealthMonitor {
	return &RedisHealthMonitor{
		rateLimitClient:  rateLimitClient,
		fallbackClient:   fallbackClient,
		fallbackSWClient: fallbackSWClient,
		retryInterval:    retryInterval,
		stopCh:           make(chan struct{}),
	}
}

// Start begins monitoring Redis health
func (h *RedisHealthMonitor) Start() {
	log.Info().Msgf("Starting Redis health monitor (retry interval: %s)", h.retryInterval)
	go h.monitorHealth()
}

// Stop stops the health monitor
func (h *RedisHealthMonitor) Stop() {
	close(h.stopCh)
}

// monitorHealth periodically checks Redis connection
func (h *RedisHealthMonitor) monitorHealth() {
	ticker := time.NewTicker(h.retryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.checkAndRestore()
		case <-h.stopCh:
			log.Info().Msg("Redis health monitor stopped")
			return
		}
	}
}

// checkAndRestore checks if Redis is available and restores connection if needed
func (h *RedisHealthMonitor) checkAndRestore() {
	// Only check if we're currently using fallback
	if h.fallbackClient != nil && !h.fallbackClient.IsRedisAvailable() {
		if h.pingRedis() {
			log.Info().Msg("Redis connection restored successfully")

			if h.fallbackClient != nil {
				h.fallbackClient.RestoreRedis()
			}

			if h.fallbackSWClient != nil {
				h.fallbackSWClient.RestoreRedis()
			}
		}
	}
}

// pingRedis attempts to ping Redis to check availability
func (h *RedisHealthMonitor) pingRedis() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.rateLimitClient.Ping(ctx).Result()
	if err != nil {
		log.Debug().Err(err).Msg("Redis ping failed during health check")
		return false
	}

	if result == "" || result != "PONG" {
		log.Debug().Msgf("Redis ping returned unexpected result: %s", result)
		return false
	}

	return true
}
