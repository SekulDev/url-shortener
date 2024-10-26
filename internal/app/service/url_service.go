package service

import (
	"url-shortener/internal/app/usecase"
	"url-shortener/internal/domain/entity"
)

type UrlService interface {
	ResolveShortUrl(url string) (*entity.Url, error)
	AddUrl(longUrl string) (*entity.Url, error)
}

type UrlServiceImpl struct {
	UrlService
	usecase     *usecase.UrlUsecaseImpl
	hashUsecase *usecase.HashUsecase
}

func NewUrlService(u *usecase.UrlUsecaseImpl, h *usecase.HashUsecase) *UrlServiceImpl {
	return &UrlServiceImpl{
		usecase:     u,
		hashUsecase: h,
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

func (u *UrlServiceImpl) AddUrl(longUrl string) (*entity.Url, error) {
	shortUrl := u.hashUsecase.GenerateHash()

	url := entity.Url{
		LongUrl: longUrl,
		ShortId: shortUrl,
	}

	err := u.usecase.AddUrlToCache(&url)
	if err != nil {
		return nil, err
	}

	return &url, nil
}
