/*
	Important - This strategy is not ready to be used by Rate Shield yet. Please ignore this if you are using Rate Shield as library.
*/

package limiter

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Request struct {
	ClientIP string
	Endpoint string
}

type LeakyBucket struct {
	Reqs          []Request
	Capacity      int
	EmptyRate     time.Duration
	StopRefilling chan struct{}
	Mutex         sync.Mutex
}

func NewLeakyBucket(capacity int, emptyRate time.Duration) *LeakyBucket {
	b := &LeakyBucket{
		Capacity:      capacity,
		EmptyRate:     emptyRate,
		StopRefilling: make(chan struct{}),
	}

	go b.processRequests()
	return b
}

func (lb *LeakyBucket) processRequests() {
	ticker := time.NewTicker(lb.EmptyRate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			lb.Mutex.Lock()
			if len(lb.Reqs) > 0 {
				log.Info().Msgf("Request Served: %d", len(lb.Reqs))
				lb.Reqs = nil
			}
			lb.Mutex.Unlock()
		case <-lb.StopRefilling:
			return
		}
	}
}

func (lb *LeakyBucket) AddRequestsToBucket(req Request) bool {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()

	if len(lb.Reqs) < lb.Capacity {
		lb.Reqs = append(lb.Reqs, req)
		return true
	}

	return false
}

func (lb *LeakyBucket) StopRefillingBucket() {
	close(lb.StopRefilling)
}
