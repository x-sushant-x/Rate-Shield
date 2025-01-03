package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/x-sushant-x/RateShield/models"
	"github.com/x-sushant-x/RateShield/utils"
)

var (
	ctx = context.Background()
)

type SlidingWindowService struct {
	redisClient *redis.ClusterClient
}

func NewSlidingWindowService(redisClient *redis.ClusterClient) SlidingWindowService {
	return SlidingWindowService{
		redisClient: redisClient,
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

	pipe := s.redisClient.TxPipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", then)

	countCmd := pipe.ZCount(ctx, key, then, fmt.Sprintf("%d", now))

	_, err := pipe.Exec(ctx)
	if err != nil {
		return -1, err
	}

	count, err := countCmd.Result()
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (s *SlidingWindowService) updateWindow(key string, now int64, windowSize time.Duration) error {
	pipe := s.redisClient.TxPipeline()

	pipe.ZAdd(ctx, key, redis.Z{
		Member: now,
		Score:  float64(now),
	})

	pipe.Expire(ctx, key, windowSize)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
