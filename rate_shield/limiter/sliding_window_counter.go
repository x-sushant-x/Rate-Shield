package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/x-sushant-x/RateShield/models"
	redisClient "github.com/x-sushant-x/RateShield/redis"
	"github.com/x-sushant-x/RateShield/utils"
)

var (
	ctx = context.Background()
)

type SlidingWindowService struct {
	client redisClient.SlidingWindowClient
}

func NewSlidingWindowService(client redisClient.SlidingWindowClient) SlidingWindowService {
	return SlidingWindowService{
		client: client,
	}
}

func (s *SlidingWindowService) processRequest(ip, endpoint string, rule *models.Rule) *models.RateLimitResponse {
	key := ip + ":" + endpoint

	now := time.Now().Unix()
	windowSize := time.Duration(rule.SlidingWindowCounterRule.WindowSize) * time.Second

	count, err := s.removeOldRequestsAndCountActiveRequests(key, now, windowSize)
	if err != nil {
		return utils.BuildRateLimitErrorResponse(500)
	}

	if count > rule.SlidingWindowCounterRule.MaxRequests {
		return utils.BuildRateLimitErrorResponse(429)
	}

	err = s.updateWindow(key, now, windowSize)
	if err != nil {
		return utils.BuildRateLimitErrorResponse(500)
	}

	return utils.BuildRateLimitSuccessResponse(rule.SlidingWindowCounterRule.MaxRequests, rule.SlidingWindowCounterRule.MaxRequests-count)
}

func (s *SlidingWindowService) removeOldRequestsAndCountActiveRequests(key string, now int64, windowSize time.Duration) (int64, error) {
	then := fmt.Sprintf("%d", now-int64(windowSize.Seconds()))

	err := s.client.ZRemRangeByScore(ctx, key, "0", then)
	if err != nil {
		return -1, err
	}

	count, err := s.client.ZCount(ctx, key, then, fmt.Sprintf("%d", now))
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (s *SlidingWindowService) updateWindow(key string, now int64, windowSize time.Duration) error {
	err := s.client.ZAdd(ctx, key, now)
	if err != nil {
		return err
	}

	err = s.client.Expire(ctx, key, windowSize)
	if err != nil {
		return err
	}

	return nil
}
