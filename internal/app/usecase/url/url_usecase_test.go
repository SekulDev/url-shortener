package usecase

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
	"url-shortener/internal/domain/entity"
	repository "url-shortener/internal/domain/repository/url"
	database "url-shortener/internal/infrastructure/database/redis"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestGetUrlFromCache_CacheHit(t *testing.T) {
	mockRedis := new(database.MockRedis)
	mockRepo := new(repository.MockUrlRepository)

	urlUsecase := NewUrlUsecase(mockRepo, mockRedis)
	url := &entity.Url{
		LongUrl: "https://example.com",
		ShortId: "abc123",
	}

	data, _ := json.Marshal(url)
	mockRedis.On("Get", "url_abc123").Return(string(data), nil)

	result, err := urlUsecase.GetUrlFromCache("abc123")

	assert.NoError(t, err)
	assert.Equal(t, url, result)
	mockRedis.AssertExpectations(t)
}

func TestGetUrlFromCache_CacheMiss(t *testing.T) {
	mockRedis := new(database.MockRedis)
	mockRepo := new(repository.MockUrlRepository)

	urlUsecase := NewUrlUsecase(mockRepo, mockRedis)
	mockRedis.On("Get", "url_abc123").Return("", redis.Nil)

	result, err := urlUsecase.GetUrlFromCache("abc123")

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockRedis.AssertExpectations(t)
}

func TestGetUrlFromCache_InvalidFormat(t *testing.T) {
	mockRedis := new(database.MockRedis)
	mockRepo := new(repository.MockUrlRepository)

	urlUsecase := NewUrlUsecase(mockRepo, mockRedis)
	mockRedis.On("Get", "url_abc123").Return("invalid-json", nil)

	result, err := urlUsecase.GetUrlFromCache("abc123")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid url format")
	mockRedis.AssertExpectations(t)
}

func TestGetUrlFromDatabase_Success(t *testing.T) {
	mockRedis := new(database.MockRedis)
	mockRepo := new(repository.MockUrlRepository)

	urlUsecase := NewUrlUsecase(mockRepo, mockRedis)
	url := &entity.Url{
		LongUrl: "https://example.com",
		ShortId: "abc123",
	}

	mockRepo.On("GetByShortUrl", "abc123").Return(url, nil)

	result, err := urlUsecase.GetUrlFromDatabase("abc123")

	assert.NoError(t, err)
	assert.Equal(t, url, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUrlFromDatabase_NotFound(t *testing.T) {
	mockRedis := new(database.MockRedis)
	mockRepo := new(repository.MockUrlRepository)

	urlUsecase := NewUrlUsecase(mockRepo, mockRedis)
	mockRepo.On("GetByShortUrl", "abc123").Return(nil, errors.New("url not found"))

	result, err := urlUsecase.GetUrlFromDatabase("abc123")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestAddUrlToCache_Success(t *testing.T) {
	mockRedis := new(database.MockRedis)
	mockRepo := new(repository.MockUrlRepository)

	urlUsecase := NewUrlUsecase(mockRepo, mockRedis)
	url := &entity.Url{
		LongUrl: "https://example.com",
		ShortId: "abc123",
	}

	data, _ := json.Marshal(url)
	mockRedis.On("Set", "url_abc123", data, time.Minute*15).Return(nil)

	err := urlUsecase.AddUrlToCache(url)

	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)
}

func TestAddUrlToDatabase_Success(t *testing.T) {
	mockRedis := new(database.MockRedis)
	mockRepo := new(repository.MockUrlRepository)

	urlUsecase := NewUrlUsecase(mockRepo, mockRedis)
	url := &entity.Url{
		LongUrl: "https://example.com",
		ShortId: "abc123",
	}

	mockRepo.On("Create", url).Return(primitive.NewObjectID(), nil)

	err := urlUsecase.AddUrlToDatabase(url)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAddUrlToDatabase_Failure(t *testing.T) {
	mockRedis := new(database.MockRedis)
	mockRepo := new(repository.MockUrlRepository)

	urlUsecase := NewUrlUsecase(mockRepo, mockRedis)
	url := &entity.Url{
		LongUrl: "https://example.com",
		ShortId: "abc123",
	}

	mockRepo.On("Create", url).Return(primitive.NilObjectID, errors.New("database error"))

	err := urlUsecase.AddUrlToDatabase(url)

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	mockRepo.AssertExpectations(t)
}
