package usecase

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	database "url-shortener/internal/infrastructure/database/redis"
)

func TestIsAllowed_Allowed(t *testing.T) {
	mockRedis := new(database.MockRedis)
	mockRedis.On("Exists", []string{"rate_limit_192.168.1.1"}).Return(int64(0), nil)

	rateLimiter := NewRateLimitUsecase(mockRedis)

	allowed := rateLimiter.IsAllowed("192.168.1.1")

	assert.True(t, allowed)
	mockRedis.AssertExpectations(t)
}

func TestIsAllowed_Disallowed(t *testing.T) {
	mockRedis := new(database.MockRedis)
	mockRedis.On("Exists", []string{"rate_limit_192.168.1.1"}).Return(int64(1), nil)

	rateLimiter := NewRateLimitUsecase(mockRedis)

	allowed := rateLimiter.IsAllowed("192.168.1.1")

	assert.False(t, allowed)
	mockRedis.AssertExpectations(t)
}

func TestDisallow_Success(t *testing.T) {
	mockRedis := new(database.MockRedis)
	disallowTime := time.Minute * 10

	mockRedis.On("Set", "rate_limit_192.168.1.1", 1, disallowTime).Return(nil)

	rateLimiter := NewRateLimitUsecase(mockRedis)

	err := rateLimiter.Disallow("192.168.1.1")

	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)
}

func TestDisallow_Failure(t *testing.T) {
	mockRedis := new(database.MockRedis)
	disallowTime := time.Minute * 10

	mockRedis.On("Set", "rate_limit_192.168.1.1", 1, disallowTime).Return(assert.AnError)

	rateLimiter := NewRateLimitUsecase(mockRedis)

	err := rateLimiter.Disallow("192.168.1.1")

	assert.Error(t, err)
	mockRedis.AssertExpectations(t)
}
