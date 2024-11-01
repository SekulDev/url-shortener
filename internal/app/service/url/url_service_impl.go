package service

import (
	hashU "url-shortener/internal/app/usecase/hash"
	ratelimitU "url-shortener/internal/app/usecase/ratelimit"
	recaptchaU "url-shortener/internal/app/usecase/recaptcha"
	urlU "url-shortener/internal/app/usecase/url"
	"url-shortener/internal/domain/apperrors"
	"url-shortener/internal/domain/entity"
	"url-shortener/pkg"
)

type urlServiceImpl struct {
	usecase          urlU.UrlUsecase
	hashUsecase      hashU.HashUsecase
	rateLimitUsecase ratelimitU.RatelimitUsecase
	recaptchaUsecase recaptchaU.RecaptchaUsecase
}

func NewUrlService(u urlU.UrlUsecase, h hashU.HashUsecase, rl ratelimitU.RatelimitUsecase, r recaptchaU.RecaptchaUsecase) UrlService {
	return &urlServiceImpl{
		usecase:          u,
		hashUsecase:      h,
		rateLimitUsecase: rl,
		recaptchaUsecase: r,
	}
}

func (u *urlServiceImpl) ResolveShortUrl(url string) (*entity.Url, error) {
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

func (u *urlServiceImpl) AddUrl(longUrl string, ip string, recaptchaToken string) (*entity.Url, error) {
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
