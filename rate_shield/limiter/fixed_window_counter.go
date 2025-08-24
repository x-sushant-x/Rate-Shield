package limiter

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
	redisClient "github.com/x-sushant-x/RateShield/redis"
	"github.com/x-sushant-x/RateShield/utils"
)

type FixedWindowService struct {
	redisClient redisClient.RedisRateLimiterClient
}

func NewFixedWindowService(client redisClient.RedisRateLimiterClient) FixedWindowService {
	return FixedWindowService{
		redisClient: client,
	}
}

func (fw *FixedWindowService) processRequest(ip, endpoint string, rule *models.Rule) *models.RateLimitResponse {
	key := fw.parseToKey(ip, endpoint)

	fixedWindow, response := fw.prepareFixedWindow(ip, endpoint, rule)
	if response != nil {
		return response
	}

	currTime := time.Now().Unix()
	return fw.handleRateLimit(fixedWindow, key, currTime)
}

func (fw *FixedWindowService) prepareFixedWindow(ip, endpoint string, rule *models.Rule) (*models.FixedWindowCounter, *models.RateLimitResponse) {
	key := fw.parseToKey(ip, endpoint)
	fixedWindow, found, err := fw.getFixedWindowFromRedis(key)
	if err != nil {
		return nil, fw.handleError(err, "error while getting fixed window")
	}

	if !found {
		fixedWindow, err := fw.spawnNewFixedWindow(ip, endpoint, rule)
		if err != nil {
			return nil, fw.handleError(err, "unable to get newly spawned fixed window from redis")
		}
		return fixedWindow, utils.BuildRateLimitSuccessResponse(fixedWindow.MaxRequests, fixedWindow.MaxRequests-1)
	}

	return fixedWindow, nil
}

func (fw *FixedWindowService) handleRateLimit(fixedWindow *models.FixedWindowCounter, key string, currTime int64) *models.RateLimitResponse {
	if currTime-fixedWindow.CreatedAt < int64(fixedWindow.Window) {
		return fw.processWithinTimeWindow(fixedWindow, key, currTime)
	}
	return fw.ResetWindow(key, currTime, fixedWindow)
}

func (fw *FixedWindowService) processWithinTimeWindow(fixedWindow *models.FixedWindowCounter, key string, currTime int64) *models.RateLimitResponse {
	if fixedWindow.CurrRequests < fixedWindow.MaxRequests {
		fixedWindow.CurrRequests++
		fixedWindow.LastAccessTime = currTime
		return fw.saveFixedWindow(key, fixedWindow)
	}
	return utils.BuildRateLimitErrorResponse(429)
}

func (fw *FixedWindowService) saveFixedWindow(key string, fixedWindow *models.FixedWindowCounter) *models.RateLimitResponse {
	err := fw.save(key, fixedWindow)
	if err != nil {
		return fw.handleError(err, "error while saving fixed window")
	}
	return utils.BuildRateLimitSuccessResponse(fixedWindow.MaxRequests, fixedWindow.MaxRequests-fixedWindow.CurrRequests)
}

func (fw *FixedWindowService) handleError(err error, msg string) *models.RateLimitResponse {
	log.Err(err).Msg(msg)
	return utils.BuildRateLimitErrorResponse(500)
}

func (fw *FixedWindowService) ResetWindow(key string, currTime int64, fixedWindow *models.FixedWindowCounter) *models.RateLimitResponse {
	fixedWindow.CurrRequests = 1
	fixedWindow.LastAccessTime = currTime
	return fw.saveFixedWindow(key, fixedWindow)
}

func (fw *FixedWindowService) getFixedWindowFromRedis(key string) (*models.FixedWindowCounter, bool, error) {
	fixedWindowFromRedis, found, err := fw.redisClient.JSONGet(key)

	if err != nil {
		log.Error().Err(err).Msg("Error fetching fixed window from Redis")
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	fixedWindow, err := utils.Unmarshal[models.FixedWindowCounter]([]byte(fixedWindowFromRedis))
	if err != nil {
		log.Err(err).Msg(err.Error())
		return nil, true, err
	}

	return &fixedWindow, true, nil
}

func (fw *FixedWindowService) makeFixedWindowCounter(ip, endpoint string, rule *models.Rule) models.FixedWindowCounter {
	return models.FixedWindowCounter{
		Endpoint:       endpoint,
		ClientIP:       ip,
		CreatedAt:      time.Now().Unix(),
		MaxRequests:    rule.FixedWindowCounterRule.MaxRequests,
		CurrRequests:   1,
		Window:         rule.FixedWindowCounterRule.Window,
		LastAccessTime: time.Now().Unix(),
	}
}

func (fw *FixedWindowService) spawnNewFixedWindow(ip, endpoint string, rule *models.Rule) (*models.FixedWindowCounter, error) {
	key := fw.parseToKey(ip, endpoint)
	fixedWindow := fw.makeFixedWindowCounter(ip, endpoint, rule)

	if err := fw.save(key, &fixedWindow); err != nil {
		log.Err(err).Msg("unable to save fixed window to redis")
		return nil, err
	}

	err := fw.redisClient.Expire(key, time.Duration(fixedWindow.Window)*time.Second)
	if err != nil {
		log.Err(err).Msg("unable to set expire time of fixed window in redis")
		return nil, err
	}
	return &fixedWindow, nil
}

func (fw *FixedWindowService) save(key string, fixedWindow *models.FixedWindowCounter) error {
	return fw.redisClient.JSONSet(key, fixedWindow)
}

func (fw *FixedWindowService) parseToKey(ip, endpoint string) string {
	return "fixed_window_" + ip + ":" + endpoint
}
