package limiter

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
	redisClient "github.com/x-sushant-x/RateShield/redis"
	"github.com/x-sushant-x/RateShield/utils"
)

const (
	TokenAddTime = time.Minute * 1
)

type Limiter struct {
	tokenBucket *TokenBucketService
	fixedWindow *FixedWindowService
}

func NewRateLimiterService(tokenBucket *TokenBucketService, fixedWindow *FixedWindowService) Limiter {
	return Limiter{
		tokenBucket: tokenBucket,
		fixedWindow: fixedWindow,
	}
}

func (l *Limiter) CheckLimit(ip, endpoint string) *models.RateLimitResponse {
	key := ip + ":" + endpoint
	rule, found, err := l.GetRule(endpoint)

	if err == nil && found {
		switch rule.Strategy {
		case "TOKEN BUCKET":
			return l.tokenBucket.processRequest(key, *rule)
		case "FIXED WINDOW COUNTER":
			return l.fixedWindow.processRequest(ip, endpoint, *rule)
		}
	}

	if err != nil {
		log.Err(err).Msg("unable to check limit")
		// Notify on slack
		return utils.BuildRateLimitSuccessResponse(0, 0)
	}

	return utils.BuildRateLimitSuccessResponse(0, 0)
}

func (l *Limiter) GetRule(key string) (*models.Rule, bool, error) {
	return redisClient.GetRule(key)
}

func (l *Limiter) StartRateLimiter() {
	log.Info().Msg("Starting Limiter Service ✅")
	l.startAddTokenJob()
}

func (l *Limiter) startAddTokenJob() {
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
