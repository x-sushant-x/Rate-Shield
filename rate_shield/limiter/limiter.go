package limiter

import (
	"net/http"
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
	tokenBucket          *TokenBucketService
	fixedWindow          *FixedWindowService
	redisRuleSvc         service.RulesService
	errorNotificationSvc *service.ErrorNotificationSVC
}

func NewRateLimiterService(tokenBucket *TokenBucketService, fixedWindow *FixedWindowService, redisRuleSvc service.RulesService) Limiter {
	return Limiter{
		tokenBucket:  tokenBucket,
		fixedWindow:  fixedWindow,
		redisRuleSvc: redisRuleSvc,
	}
}

func (l *Limiter) CheckLimit(ip, endpoint string) *models.RateLimitResponse {
	key := ip + ":" + endpoint
	rule, found, err := l.GetRule(endpoint)

	if err == nil && found {
		switch rule.Strategy {
		case "TOKEN BUCKET":
			return l.processTokenBucketReq(key, rule)
		case "FIXED WINDOW COUNTER":
			return l.processFixedWindowReq(ip, endpoint, rule)
		}
	}

	if err != nil {
		log.Err(err).Msg("error while getting rule for api endpoint: " + endpoint)
		l.errorNotificationSvc.SendErrorNotification(err.Error(), time.Now(), ip, endpoint, *rule)
		return utils.BuildRateLimitErrorResponse(http.StatusInternalServerError)
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
	l.tokenBucket.startAddTokenJob()
}
