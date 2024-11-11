package limiter

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
	"github.com/x-sushant-x/RateShield/service"
	"github.com/x-sushant-x/RateShield/utils"
)

const (
	TokenAddTime = time.Minute * 1
)

type Limiter struct {
	tokenBucket   *TokenBucketService
	fixedWindow   *FixedWindowService
	slidingWindow *SlidingWindowService
	redisRuleSvc  service.RulesService
	cachedRules   map[string]*models.Rule
}

func NewRateLimiterService(
	tokenBucket *TokenBucketService, fixedWindow *FixedWindowService, slidingWindow *SlidingWindowService, redisRuleSvc service.RulesService) Limiter {

	return Limiter{
		tokenBucket:   tokenBucket,
		fixedWindow:   fixedWindow,
		redisRuleSvc:  redisRuleSvc,
		slidingWindow: slidingWindow,
		// This is initialized later in StartRateLimiter() function
		cachedRules: nil,
	}
}

func (l *Limiter) CheckLimit(ip, endpoint string) *models.RateLimitResponse {
	key := ip + ":" + endpoint
	rule, found := l.cachedRules[endpoint]

	if found {
		switch rule.Strategy {
		case "TOKEN BUCKET":
			return l.processTokenBucketReq(key, rule)
		case "FIXED WINDOW COUNTER":
			return l.processFixedWindowReq(ip, endpoint, rule)
		case "SLIDING WINDOW COUNTER":
			return l.processSlidingWindowReq(ip, endpoint, rule)
		}
	}

	if !found {
		return utils.BuildRateLimitSuccessResponse(0, 0)
	}

	return utils.BuildRateLimitSuccessResponse(0, 0)
}

func (l *Limiter) processTokenBucketReq(key string, rule *models.Rule) *models.RateLimitResponse {
	resp := l.tokenBucket.processRequest(key, rule)

	if resp.Success {
		return resp
	}

	if rule.AllowOnError {
		return utils.BuildRateLimitSuccessResponse(0, 0)
	}

	return resp
}

func (l *Limiter) processFixedWindowReq(ip, endpoint string, rule *models.Rule) *models.RateLimitResponse {
	resp := l.fixedWindow.processRequest(ip, endpoint, rule)

	if resp.Success {
		return resp
	}

	if rule.AllowOnError {
		return utils.BuildRateLimitSuccessResponse(0, 0)
	}

	return resp
}

func (l *Limiter) processSlidingWindowReq(ip, endpoint string, rule *models.Rule) *models.RateLimitResponse {
	resp := l.slidingWindow.processRequest(ip, endpoint, rule)

	if resp.Success {
		return resp
	}

	if rule.AllowOnError {
		return utils.BuildRateLimitSuccessResponse(0, 0)
	}

	return resp
}

func (l *Limiter) GetRule(key string) (*models.Rule, bool, error) {
	return l.redisRuleSvc.GetRule(key)
}

func (l *Limiter) StartRateLimiter() {
	log.Info().Msg("Starting limiter service ✅")
	l.cachedRules = l.redisRuleSvc.CacheRulesLocally()
	l.tokenBucket.startAddTokenJob()
	go l.listenToRulesUpdate()

}

func (l *Limiter) listenToRulesUpdate() {
	updatesChannel := make(chan string)
	go l.redisRuleSvc.ListenToRulesUpdate(updatesChannel)

	for {
		data := <-updatesChannel

		if data == "UpdateRules" {
			l.cachedRules = l.redisRuleSvc.CacheRulesLocally()
			log.Info().Msg("Rules Updated Successfully")
		}
	}
}
