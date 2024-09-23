package limiter

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
	redisClient "github.com/x-sushant-x/RateShield/redis"
	"github.com/x-sushant-x/RateShield/utils"
)

type FixedWindowService struct{}

func NewFixedWindowService() FixedWindowService {
	return FixedWindowService{}
}

func (fw *FixedWindowService) processRequest(ip, endpoint string, rule models.Rule) *models.RateLimitResponse {
	key := fw.parseToKey(ip, endpoint)

	fixedWindow, found, err := fw.getFixedWindowFromRedis(key)
	if err != nil {
		log.Err(err).Msgf("error while getting fixed window")
		return utils.BuildRateLimitErrorResponse(500)
	}

	if !found {
		fixedWindow, err := fw.spawnNewFixedWindow(ip, endpoint, rule)
		if err != nil {
			log.Err(err).Msg("unable to get newly spawned fixed window from redis")
			return utils.BuildRateLimitErrorResponse(500)
		}
		return utils.BuildRateLimitSuccessResponse(fixedWindow.MaxRequests, fixedWindow.MaxRequests-1)
	}

	currTime := time.Now().Unix()

	if currTime-fixedWindow.LastAccessTime < int64(fixedWindow.Window) {
		if fixedWindow.CurrRequests < fixedWindow.MaxRequests {
			fixedWindow.CurrRequests++
			fixedWindow.LastAccessTime = currTime
			err := fw.save(key, fixedWindow)
			if err != nil {
				log.Err(err).Msg("error while saving fixed window")
				return utils.BuildRateLimitErrorResponse(500)
			}
			return utils.BuildRateLimitSuccessResponse(fixedWindow.MaxRequests, fixedWindow.MaxRequests-fixedWindow.CurrRequests)
		}
		return utils.BuildRateLimitErrorResponse(429)

	} else {
		fixedWindow.CurrRequests = 1
		fixedWindow.LastAccessTime = currTime

		err := fw.save(key, fixedWindow)
		if err != nil {
			log.Err(err).Msg("error while saving fixed window")
			return utils.BuildRateLimitErrorResponse(500)

		}
		return utils.BuildRateLimitSuccessResponse(fixedWindow.MaxRequests, fixedWindow.MaxRequests-fixedWindow.CurrRequests)
	}
}

func (fw *FixedWindowService) ResetWindow(key string, currTime int64, window *models.FixedWindowCounter) *models.RateLimitResponse {
	window.CurrRequests = 1
	window.LastAccessTime = currTime

	err := fw.save(key, window)
	if err != nil {
		log.Err(err).Msg("error while saving fixed window")
		return utils.BuildRateLimitErrorResponse(500)

	}
	return utils.BuildRateLimitSuccessResponse(window.MaxRequests, window.MaxRequests-window.CurrRequests)

}

func (fw *FixedWindowService) getFixedWindowFromRedis(key string) (*models.FixedWindowCounter, bool, error) {
	data, found, err := redisClient.GetFixedWindowJSONObject(key)

	if err != nil {
		log.Error().Err(err).Msg("Error fetching fixed window from Redis")
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	fixedWindow, err := unmarshalFixedWindow(data)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling fixed window from Redis")
		return nil, false, err
	}
	return fixedWindow, true, nil
}

func (fw *FixedWindowService) spawnNewFixedWindow(ip, endpoint string, rule models.Rule) (*models.FixedWindowCounter, error) {
	key := fw.parseToKey(ip, endpoint)
	fixedWindow := models.FixedWindowCounter{
		Endpoint:       endpoint,
		ClientIP:       ip,
		CreatedAt:      time.Now().Unix(),
		MaxRequests:    rule.FixedWindowCounterRule.MaxRequests,
		CurrRequests:   1,
		Window:         rule.FixedWindowCounterRule.Window,
		LastAccessTime: time.Now().Unix(),
	}

	if err := fw.save(key, &fixedWindow); err != nil {
		log.Err(err).Msg("unable to save fixed window to redis")
		return nil, err
	}

	err := redisClient.SetFixedWindowExpireTime(key, time.Duration(fixedWindow.Window)*time.Second)
	if err != nil {
		log.Err(err).Msg("unable to set expire time of fixed window in redis")
		return nil, err
	}
	return &fixedWindow, nil
}

func (fw *FixedWindowService) save(key string, fixedWindow *models.FixedWindowCounter) error {
	return redisClient.SetFixedWindowJSONObject(key, fixedWindow)
}

func unmarshalFixedWindow(data []byte) (*models.FixedWindowCounter, error) {
	fixedWindow := models.FixedWindowCounter{}

	if err := json.Unmarshal(data, &fixedWindow); err != nil {
		log.Error().Err(err).Msg("Error unmarshalling fixed window data")
		return &fixedWindow, err
	}
	return &fixedWindow, nil
}

func (fw *FixedWindowService) parseToKey(ip, endpoint string) string {
	return ip + ":" + endpoint
}
