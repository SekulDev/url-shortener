package service

import (
	"url-shortener/internal/app/usecase"
	"url-shortener/internal/domain/apperrors"
	"url-shortener/internal/domain/entity"
	"url-shortener/pkg"
)

type UrlService interface {
	ResolveShortUrl(url string) (*entity.Url, error)
	AddUrl(longUrl string, ip string, recaptchaToken string) (*entity.Url, error)
}

type UrlServiceImpl struct {
	UrlService
	usecase          *usecase.UrlUsecaseImpl
	hashUsecase      *usecase.HashUsecase
	rateLimitUsecase *usecase.RateLimitUsecaseImpl
	recaptchaUsecase *usecase.RecaptchaUsecaseImpl
}

func NewUrlService(u *usecase.UrlUsecaseImpl, h *usecase.HashUsecase, rl *usecase.RateLimitUsecaseImpl, r *usecase.RecaptchaUsecaseImpl) *UrlServiceImpl {
	return &UrlServiceImpl{
		usecase:          u,
		hashUsecase:      h,
		rateLimitUsecase: rl,
		recaptchaUsecase: r,
	}
}

func (u *UrlServiceImpl) ResolveShortUrl(url string) (*entity.Url, error) {
	cacheUrl, cacheErr := u.usecase.GetUrlFromCache(url)
	if cacheErr != nil {
		return nil, cacheErr
	}

	if cacheUrl != nil {
		return cacheUrl, nil
	}

	dbUrl, dbErr := u.usecase.GetUrlFromDatabase(url)
	if dbErr != nil {
		return nil, dbErr
	}

	_ = u.usecase.AddUrlToCache(dbUrl)

	return dbUrl, nil
}

func (u *UrlServiceImpl) AddUrl(longUrl string, ip string, recaptchaToken string) (*entity.Url, error) {
	recaptchaErr := u.recaptchaUsecase.Verify(recaptchaToken)
	if recaptchaErr != nil {
		return nil, apperrors.RecaptchaError
	}

	if !pkg.IsValidUrl(longUrl) {
		return nil, apperrors.InvalidUrlError
	}

	if !u.rateLimitUsecase.IsAllowed(ip) {
		return nil, apperrors.TooManyRequests
	}

	shortUrl := u.hashUsecase.GenerateHash()

	url := entity.Url{
		LongUrl: longUrl,
		ShortId: shortUrl,
	}

	err := u.usecase.AddUrlToDatabase(&url)
	if err != nil {
		return nil, err
	}

	_ = u.rateLimitUsecase.Disallow(ip)

	return &url, nil
}
