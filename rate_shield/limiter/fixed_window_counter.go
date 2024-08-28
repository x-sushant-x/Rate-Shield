package limiter

import (
	"net/http"
	"sync"
	"time"
)

type FixedWindowService struct {
	requests      int
	window        time.Duration
	visitorRecord map[string]*visitor
	mu            sync.Mutex
}

type visitor struct {
	lastAccessTime time.Time
	count          int
}

func NewFixedWindowService(requests int, window time.Duration) FixedWindowService {
	return FixedWindowService{
		requests:      requests,
		window:        window,
		visitorRecord: make(map[string]*visitor),
	}
}

func (fw *FixedWindowService) processRequest(ip, endpoint string) int {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()

	key := ip + ":" + endpoint

	if v, exists := fw.visitorRecord[key]; exists {
		if now.Sub(v.lastAccessTime) > fw.window {
			v.count = 1
			v.lastAccessTime = now
		} else if v.count < fw.requests {
			v.count++
		} else {
			return http.StatusTooManyRequests
		}
	} else {
		fw.visitorRecord[key] = &visitor{lastAccessTime: time.Now(), count: 1}
	}
	return http.StatusOK
}
