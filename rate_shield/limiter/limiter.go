package limiter

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

var (
	RateLimiter *Limiter
)

type Limiter struct {
	TokenBucket TokenBucketService
}

func StartSvc() {
	RateLimiter = &Limiter{}
	StartAddTokenJob()
}

func (l *Limiter) CheckLimit(ip, endpoint string) int {
	return l.TokenBucket.ProcessRequest(ip, endpoint)
}

func StartAddTokenJob() {
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
}
