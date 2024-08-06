package limiter

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

var (
	RateLimiter *Limiter
)

type Limiter struct {
	TokenBucket TokenBuckets
}

func StartSvc() {
	RateLimiter = &Limiter{
		TokenBucket: *getTokenBucketInstance(),
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	_, err = s.NewJob(gocron.DurationJob(time.Second*10), gocron.NewTask(func() {
		RateLimiter.TokenBucket.AddTokens()
	}))

	if err != nil {
		panic(err)
	}

	s.Start()
	fmt.Println("Add token job started")
}

func (l *Limiter) CheckLimit(ip string) bool {
	return l.TokenBucket.ProcessRequest(ip)
}
