package fallback

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// RulesRefresher interface for components that can refresh cached rules
type RulesRefresher interface {
	RefreshCachedRules()
}

// RulesHealthMonitor monitors Redis rules instance health and manages reconnection
type RulesHealthMonitor struct {
	rulesClient    *redis.Client
	rulesRefresher RulesRefresher
	retryInterval  time.Duration
	stopCh         chan struct{}
	isAvailable    bool
}

// NewRulesHealthMonitor creates a new rules health monitor
func NewRulesHealthMonitor(
	rulesClient *redis.Client,
	rulesRefresher RulesRefresher,
	retryInterval time.Duration,
) *RulesHealthMonitor {
	return &RulesHealthMonitor{
		rulesClient:    rulesClient,
		rulesRefresher: rulesRefresher,
		retryInterval:  retryInterval,
		stopCh:         make(chan struct{}),
		isAvailable:    true, // Start as available
	}
}

// Start begins monitoring Redis rules instance health
func (h *RulesHealthMonitor) Start() {
	log.Info().Msgf("Starting Rules Redis health monitor (retry interval: %s)", h.retryInterval)
	go h.monitorHealth()
}

// Stop stops the health monitor
func (h *RulesHealthMonitor) Stop() {
	close(h.stopCh)
}

// monitorHealth periodically checks Redis rules instance connection
func (h *RulesHealthMonitor) monitorHealth() {
	ticker := time.NewTicker(h.retryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.checkAndRefresh()
		case <-h.stopCh:
			log.Info().Msg("Rules Redis health monitor stopped")
			return
		}
	}
}

// checkAndRefresh checks if Rules Redis is available and refreshes rules if recovered
func (h *RulesHealthMonitor) checkAndRefresh() {
	wasAvailable := h.isAvailable
	h.isAvailable = h.pingRulesRedis()

	// If it was down and is now up, refresh rules
	if !wasAvailable && h.isAvailable {
		log.Info().Msg("Rules Redis connection restored, refreshing cached rules")
		h.rulesRefresher.RefreshCachedRules()
	}

	// If it was up and is now down, log warning
	if wasAvailable && !h.isAvailable {
		log.Warn().Msg("Rules Redis connection lost, continuing with cached rules")
	}
}

// pingRulesRedis attempts to ping Rules Redis to check availability
func (h *RulesHealthMonitor) pingRulesRedis() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.rulesClient.Ping(ctx).Result()
	if err != nil {
		log.Debug().Err(err).Msg("Rules Redis ping failed during health check")
		return false
	}

	if result == "" || result != "PONG" {
		log.Debug().Msgf("Rules Redis ping returned unexpected result: %s", result)
		return false
	}

	return true
}
