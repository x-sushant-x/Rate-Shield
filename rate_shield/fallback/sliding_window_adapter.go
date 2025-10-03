package fallback

import (
	"context"
	"sort"
	"strconv"
	"sync"
	"time"
)

// SlidingWindowEntry represents a timestamp entry in the sliding window
type SlidingWindowEntry struct {
	Timestamp int64
	ExpiresAt time.Time
}

// InMemorySlidingWindowStore provides sorted set-like operations for sliding windows
type InMemorySlidingWindowStore struct {
	windows map[string][]SlidingWindowEntry
	mutex   sync.RWMutex
	stopCh  chan struct{}
}

// NewInMemorySlidingWindowStore creates a new sliding window store
func NewInMemorySlidingWindowStore() *InMemorySlidingWindowStore {
	store := &InMemorySlidingWindowStore{
		windows: make(map[string][]SlidingWindowEntry),
		stopCh:  make(chan struct{}),
	}

	// Start background cleanup
	go store.cleanupExpiredWindows()

	return store
}

// ZRemRangeByScore removes entries with scores between min and max (inclusive)
func (s *InMemorySlidingWindowStore) ZRemRangeByScore(ctx context.Context, key, min, max string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	entries, found := s.windows[key]
	if !found {
		return nil
	}

	minScore, err := strconv.ParseInt(min, 10, 64)
	if err != nil {
		return err
	}

	maxScore, err := strconv.ParseInt(max, 10, 64)
	if err != nil {
		return err
	}

	// Filter out entries in the range
	filtered := make([]SlidingWindowEntry, 0)
	for _, entry := range entries {
		if entry.Timestamp < minScore || entry.Timestamp > maxScore {
			filtered = append(filtered, entry)
		}
	}

	s.windows[key] = filtered
	return nil
}

// ZCount counts entries with scores between min and max (inclusive)
func (s *InMemorySlidingWindowStore) ZCount(ctx context.Context, key, min, max string) (int64, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	entries, found := s.windows[key]
	if !found {
		return 0, nil
	}

	minScore, err := strconv.ParseInt(min, 10, 64)
	if err != nil {
		return 0, err
	}

	maxScore, err := strconv.ParseInt(max, 10, 64)
	if err != nil {
		return 0, err
	}

	count := int64(0)
	for _, entry := range entries {
		if entry.Timestamp >= minScore && entry.Timestamp <= maxScore {
			count++
		}
	}

	return count, nil
}

// ZAdd adds an entry with the given score
func (s *InMemorySlidingWindowStore) ZAdd(ctx context.Context, key string, timestamp int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	entry := SlidingWindowEntry{
		Timestamp: timestamp,
		ExpiresAt: time.Time{}, // Will be set by Expire
	}

	s.windows[key] = append(s.windows[key], entry)

	// Keep sorted by timestamp for efficiency
	sort.Slice(s.windows[key], func(i, j int) bool {
		return s.windows[key][i].Timestamp < s.windows[key][j].Timestamp
	})

	return nil
}

// Expire sets expiration time for all entries in a key
func (s *InMemorySlidingWindowStore) Expire(ctx context.Context, key string, expiration time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	entries, found := s.windows[key]
	if !found {
		return nil
	}

	expiresAt := time.Now().Add(expiration)
	for i := range entries {
		entries[i].ExpiresAt = expiresAt
	}

	return nil
}

// cleanupExpiredWindows runs periodically to remove expired windows
func (s *InMemorySlidingWindowStore) cleanupExpiredWindows() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mutex.Lock()
			now := time.Now()
			for key, entries := range s.windows {
				// Remove if all entries are expired (check first entry since they all have same expiry)
				if len(entries) > 0 && !entries[0].ExpiresAt.IsZero() && now.After(entries[0].ExpiresAt) {
					delete(s.windows, key)
				}
			}
			s.mutex.Unlock()
		case <-s.stopCh:
			return
		}
	}
}

// Stop stops the cleanup goroutine
func (s *InMemorySlidingWindowStore) Stop() {
	close(s.stopCh)
}

// Clear removes all data
func (s *InMemorySlidingWindowStore) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.windows = make(map[string][]SlidingWindowEntry)
}
