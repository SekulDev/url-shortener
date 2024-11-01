package usecase

import (
	"fmt"
	"time"
	database "url-shortener/internal/infrastructure/database/redis"
)

type rateLimitUsecaseImpl struct {
	redis        database.Redis
	disallowTime time.Duration
}

func NewRateLimitUsecase(redis database.Redis) RatelimitUsecase {
	return &rateLimitUsecaseImpl{
		redis:        redis,
		disallowTime: time.Minute * 10,
	}
}

func (rl *rateLimitUsecaseImpl) IsAllowed(ip string) bool {
	val, err := rl.redis.Exists(fmt.Sprintf(rateLimitKey, ip))
	if err != nil {
		return false
	}
	return val == 0
}

func (rl *rateLimitUsecaseImpl) Disallow(ip string) error {
	return rl.redis.Set(fmt.Sprintf(rateLimitKey, ip), 1, rl.disallowTime)
}
