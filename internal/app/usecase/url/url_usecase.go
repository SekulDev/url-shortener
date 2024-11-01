package usecase

import (
	"url-shortener/internal/domain/entity"
)

const redisKey = "url_%s"

type UrlUsecase interface {
	GetUrlFromCache(url string) (*entity.Url, error)
	GetUrlFromDatabase(url string) (*entity.Url, error)
	AddUrlToCache(url *entity.Url) error
	AddUrlToDatabase(url *entity.Url) error
}
