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
	tokenBucket  *TokenBucketService
	fixedWindow  *FixedWindowService
	redisRuleSvc service.RulesService
	cachedRules  map[string]*models.Rule
}

func NewRateLimiterService(tokenBucket *TokenBucketService, fixedWindow *FixedWindowService, redisRuleSvc service.RulesService) Limiter {

	return Limiter{
		tokenBucket:  tokenBucket,
		fixedWindow:  fixedWindow,
		redisRuleSvc: redisRuleSvc,
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

func (l *Limiter) GetRule(key string) (*models.Rule, bool, error) {
	return l.redisRuleSvc.GetRule(key)
}

func (l *Limiter) StartRateLimiter() {
	log.Info().Msg("Starting Limiter Service ✅")
	l.cachedRules = l.cacheRulesLoally()
	l.tokenBucket.startAddTokenJob()
}

func (l *Limiter) cacheRulesLoally() map[string]*models.Rule {
	rules, err := l.redisRuleSvc.GetAllRules()
	if err != nil {
		log.Err(err).Msg("Unable to cache all rules locally")
	}

	cachedRules := make(map[string]*models.Rule)

	for _, rule := range rules {
		cachedRules[rule.APIEndpoint] = &rule
	}

	log.Info().Msg("Rules locally cached ✅")
	return cachedRules
}
