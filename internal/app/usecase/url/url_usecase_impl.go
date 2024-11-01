package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
	"url-shortener/internal/domain/entity"
	repository "url-shortener/internal/domain/repository/url"
	database "url-shortener/internal/infrastructure/database/redis"
)

type urlUsecaseImpl struct {
	repo      repository.UrlRepository
	redis     database.Redis
	cacheTime time.Duration
}

func NewUrlUsecase(r repository.UrlRepository, redis database.Redis) UrlUsecase {
	return &urlUsecaseImpl{
		repo:      r,
		redis:     redis,
		cacheTime: time.Minute * 15,
	}
}

func (u *urlUsecaseImpl) GetUrlFromCache(url string) (*entity.Url, error) {
	var result entity.Url
	data, err := u.redis.Get(fmt.Sprintf(redisKey, url))
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

func (u *urlUsecaseImpl) GetUrlFromDatabase(url string) (*entity.Url, error) {
	result, err := u.repo.GetByShortUrl(url)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("url not found")
	}

	return result, nil
}

func (u *urlUsecaseImpl) AddUrlToCache(url *entity.Url) error {
	data, jsonErr := json.Marshal(url)
	if jsonErr != nil {
		return jsonErr
	}

	err := u.redis.Set(fmt.Sprintf(redisKey, url.ShortId), data, u.cacheTime)
	return err
}

func (u *urlUsecaseImpl) AddUrlToDatabase(url *entity.Url) error {
	_, err := u.repo.Create(url)
	if err != nil {
		return err
	}

	return err
}
