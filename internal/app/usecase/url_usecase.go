package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
	"url-shortener/internal/domain/entity"
	"url-shortener/internal/domain/repository"
	database "url-shortener/internal/infrastructure/database/redis"
)

const redisKey = "url_%s"

type UrlUsecase interface {
	GetUrlFromCache(url string) (*entity.Url, error)
	GetUrlFromDatabase(url string) (*entity.Url, error)
	AddUrlToCache(url *entity.Url) error
	AddUrlToDatabase(url *entity.Url) error
}

type UrlUsecaseImpl struct {
	UrlUsecase
	repo      *repository.MongoUrlRepository
	redis     redis.Client
	cacheTime time.Duration
}

func NewUrlUsecase(r *repository.MongoUrlRepository, redis *database.RedisClient) *UrlUsecaseImpl {
	return &UrlUsecaseImpl{
		repo:      r,
		redis:     *redis.Client,
		cacheTime: time.Minute * 15,
	}
}

func (u *UrlUsecaseImpl) GetUrlFromCache(url string) (*entity.Url, error) {
	var result entity.Url
	data, err := u.redis.Get(fmt.Sprintf(redisKey, url)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		if err := json.Unmarshal([]byte(data), &result); err != nil {
			return nil, errors.New("invalid url format")
		}
		return &result, nil
	}
}

func (u *UrlUsecaseImpl) GetUrlFromDatabase(url string) (*entity.Url, error) {
	result, err := u.repo.GetByShortUrl(url)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("url not found")
	}

	return result, nil
}

func (u *UrlUsecaseImpl) AddUrlToCache(url *entity.Url) error {
	data, jsonErr := json.Marshal(url)
	if jsonErr != nil {
		return jsonErr
	}

	err := u.redis.Set(fmt.Sprintf(redisKey, url.ShortId), data, u.cacheTime).Err()
	return err
}

func (u *UrlUsecaseImpl) AddUrlToDatabase(url *entity.Url) error {
	_, err := u.repo.Create(url)
	if err != nil {
		return err
	}

	return err
}
