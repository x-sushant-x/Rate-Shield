package fallback

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// ExpiringValue holds a value with its expiration time
type ExpiringValue struct {
	Value     interface{}
	ExpiresAt time.Time
	HasExpiry bool
}

// InMemoryRateLimitStore provides in-memory storage for rate limiting data
type InMemoryRateLimitStore struct {
	data   map[string]*ExpiringValue
	mutex  sync.RWMutex
	stopCh chan struct{}
}

// NewInMemoryRateLimitStore creates a new in-memory store
func NewInMemoryRateLimitStore() *InMemoryRateLimitStore {
	store := &InMemoryRateLimitStore{
		data:   make(map[string]*ExpiringValue),
		stopCh: make(chan struct{}),
	}

	// Start background cleanup goroutine
	go store.cleanupExpiredKeys()

	return store
}

// JSONSet stores a value as JSON
func (s *InMemoryRateLimitStore) JSONSet(key string, val interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// If the value already exists and has an expiry, preserve it
	var existingExpiry time.Time
	var hasExpiry bool
	if existing, ok := s.data[key]; ok && existing.HasExpiry {
		existingExpiry = existing.ExpiresAt
		hasExpiry = true
	}

	s.data[key] = &ExpiringValue{
		Value:     val,
		ExpiresAt: existingExpiry,
		HasExpiry: hasExpiry,
	}

	return nil
}

// JSONGet retrieves a value as JSON string
func (s *InMemoryRateLimitStore) JSONGet(key string) (string, bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	val, found := s.data[key]
	if !found {
		return "", false, nil
	}

	// Check if expired
	if val.HasExpiry && time.Now().After(val.ExpiresAt) {
		return "", false, nil
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(val.Value)
	if err != nil {
		log.Error().Err(err).Msg("Error marshaling value to JSON in memory store")
		return "", false, err
	}

	return string(jsonBytes), true, nil
}

// Expire sets an expiration time on a key
func (s *InMemoryRateLimitStore) Expire(key string, expiration time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	val, found := s.data[key]
	if !found {
		// Key doesn't exist, nothing to expire
		return nil
	}

	val.ExpiresAt = time.Now().Add(expiration)
	val.HasExpiry = true

	return nil
}

// Delete removes a key from the store
func (s *InMemoryRateLimitStore) Delete(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, key)
	return nil
}

// GetKeys returns all keys matching a pattern (simplified - only supports prefix with *)
func (s *InMemoryRateLimitStore) GetKeys(pattern string) []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var keys []string

	// Simple pattern matching: if pattern ends with *, match prefix
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		for key := range s.data {
			// Check if expired
			val := s.data[key]
			if val.HasExpiry && time.Now().After(val.ExpiresAt) {
				continue
			}

			if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
				keys = append(keys, key)
			}
		}
	} else {
		// Exact match
		if val, found := s.data[pattern]; found {
			if !val.HasExpiry || time.Now().Before(val.ExpiresAt) {
				keys = append(keys, pattern)
			}
		}
	}

	return keys
}

// cleanupExpiredKeys runs periodically to remove expired keys
func (s *InMemoryRateLimitStore) cleanupExpiredKeys() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mutex.Lock()
			now := time.Now()
			for key, val := range s.data {
				if val.HasExpiry && now.After(val.ExpiresAt) {
					delete(s.data, key)
				}
			}
			s.mutex.Unlock()
		case <-s.stopCh:
			return
		}
	}
}

// Stop stops the cleanup goroutine
func (s *InMemoryRateLimitStore) Stop() {
	close(s.stopCh)
}

// Clear removes all data from the store
func (s *InMemoryRateLimitStore) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data = make(map[string]*ExpiringValue)
}
