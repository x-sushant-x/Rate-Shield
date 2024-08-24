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
}

func NewRateLimiterService(tokenBucket TokenBucketService) Limiter {
	return Limiter{
		tokenBucket: tokenBucket,
	}
}

func (l Limiter) CheckLimit(ip, endpoint string) int {
	return l.tokenBucket.ProcessRequest(ip, endpoint)
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
		l.tokenBucket.AddTokens()
	}))

	if err != nil {
		panic(err)
	}

	s.Start()
}
