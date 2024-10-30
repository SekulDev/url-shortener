package usecase

import (
	"fmt"
	"time"
	database "url-shortener/internal/infrastructure/database/redis"
)

const rateLimitKey = "rate_limit_%s"

type RateLimitUsecase interface {
	IsAllowed(ip string) bool
	Disallow(ip string) error
}

type RateLimitUsecaseImpl struct {
	redis        database.RedisClient
	disallowTime time.Duration
}

func NewRateLimitUsecase(redis *database.RedisClient) *RateLimitUsecaseImpl {
	return &RateLimitUsecaseImpl{
		redis:        *redis,
		disallowTime: time.Minute * 10,
	}
}

func (rl *RateLimitUsecaseImpl) IsAllowed(ip string) bool {
	return rl.redis.Client.Exists(fmt.Sprintf(rateLimitKey, ip)).Val() == 0
}

func (rl *RateLimitUsecaseImpl) Disallow(ip string) error {
	return rl.redis.Client.Set(fmt.Sprintf(rateLimitKey, ip), 1, rl.disallowTime).Err()
}
