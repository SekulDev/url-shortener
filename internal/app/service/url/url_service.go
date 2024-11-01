package service

import (
	"url-shortener/internal/domain/entity"
)

type UrlService interface {
	ResolveShortUrl(url string) (*entity.Url, error)
	AddUrl(longUrl string, ip string, recaptchaToken string) (*entity.Url, error)
}
