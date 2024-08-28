package limiter

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

const (
	TokenAddTime = time.Minute * 1
)

type Limiter struct {
	tokenBucket TokenBucketService
	fixedWindow *FixedWindowService
}

func NewRateLimiterService(tokenBucket TokenBucketService, fixedWindow *FixedWindowService) Limiter {
	return Limiter{
		tokenBucket: tokenBucket,
		fixedWindow: fixedWindow,
	}
}

func (l Limiter) CheckLimit(ip, endpoint string) int {
	return l.fixedWindow.processRequest(ip, endpoint)
}

func (l Limiter) StartRateLimiter() {
	l.startAddTokenJob()
}

func (l Limiter) startAddTokenJob() {
	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	_, err = s.NewJob(gocron.DurationJob(TokenAddTime), gocron.NewTask(func() {
		l.tokenBucket.addTokens()
	}))

	if err != nil {
		panic(err)
	}

	s.Start()
}
