package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/x-sushant-x/RateShield/models"
	"github.com/x-sushant-x/RateShield/utils"
)

const (
	MaxRequest = 10
	WindowSize = 60 * time.Second
)

var (
	ctx = context.Background()
)

type SlidingWindowService struct {
	redisClient *redis.Client
}

func NewSlidingWindowService(redisClient *redis.Client) SlidingWindowService {
	return SlidingWindowService{
		redisClient: redisClient,
	}
}

func (s *SlidingWindowService) ProcessRequest(key string) *models.RateLimitResponse {
	now := time.Now().Unix()

	pipe := s.redisClient.TxPipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-int64(WindowSize.Seconds())))

	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: now,
	})

	countCmd := pipe.ZCount(ctx, key, fmt.Sprintf("%d", now-int64(WindowSize.Seconds())), fmt.Sprintf("%d", now))

	pipe.Expire(ctx, key, WindowSize)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return utils.BuildRateLimitErrorResponse(500)
	}

	count, err := countCmd.Result()
	if err != nil {
		return utils.BuildRateLimitErrorResponse(500)
	}

	if count > MaxRequest {
		return utils.BuildRateLimitErrorResponse(429)
	}

	// TODO -> Change 999 from actual data when rule is configured.
	return utils.BuildRateLimitSuccessResponse(999, 999)
}
