package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	hashU "url-shortener/internal/app/usecase/hash"
	ratelimitU "url-shortener/internal/app/usecase/ratelimit"
	recaptchaU "url-shortener/internal/app/usecase/recaptcha"
	urlU "url-shortener/internal/app/usecase/url"
	"url-shortener/internal/domain/apperrors"
	"url-shortener/internal/domain/entity"
)

func TestResolveShortUrl_CacheHit(t *testing.T) {
	mockUrlUsecase := new(urlU.MockUrlUsecase)
	mockHashUsecase := new(hashU.MockHashUsecase)
	mockRatelimitUsecase := new(ratelimitU.MockRatelimitUsecase)
	mockRecaptchaUsecase := new(recaptchaU.MockRecaptchaUsecase)
	urlService := NewUrlService(mockUrlUsecase, mockHashUsecase, mockRatelimitUsecase, mockRecaptchaUsecase)

	url := "abc123"
	expectedUrl := &entity.Url{
		LongUrl: "https://example.com",
		ShortId: "abc123",
	}

	mockUrlUsecase.On("GetUrlFromCache", url).Return(expectedUrl, nil)

	result, err := urlService.ResolveShortUrl(url)

	assert.NoError(t, err)
	assert.Equal(t, expectedUrl, result)
	mockUrlUsecase.AssertExpectations(t)
}

func TestResolveShortUrl_CacheMiss_DatabaseHit(t *testing.T) {
	mockUrlUsecase := new(urlU.MockUrlUsecase)
	mockHashUsecase := new(hashU.MockHashUsecase)
	mockRatelimitUsecase := new(ratelimitU.MockRatelimitUsecase)
	mockRecaptchaUsecase := new(recaptchaU.MockRecaptchaUsecase)
	urlService := NewUrlService(mockUrlUsecase, mockHashUsecase, mockRatelimitUsecase, mockRecaptchaUsecase)

	url := "abc123"
	expectedUrl := &entity.Url{
		LongUrl: "https://example.com",
		ShortId: "abc123",
	}

	mockUrlUsecase.On("GetUrlFromCache", url).Return(nil, nil)
	mockUrlUsecase.On("GetUrlFromDatabase", url).Return(expectedUrl, nil)
	mockUrlUsecase.On("AddUrlToCache", expectedUrl).Return(nil)

	result, err := urlService.ResolveShortUrl(url)

	assert.NoError(t, err)
	assert.Equal(t, expectedUrl, result)
	mockUrlUsecase.AssertExpectations(t)
}

func TestResolveShortUrl_DatabaseMiss(t *testing.T) {
	mockUrlUsecase := new(urlU.MockUrlUsecase)
	mockHashUsecase := new(hashU.MockHashUsecase)
	mockRatelimitUsecase := new(ratelimitU.MockRatelimitUsecase)
	mockRecaptchaUsecase := new(recaptchaU.MockRecaptchaUsecase)
	urlService := NewUrlService(mockUrlUsecase, mockHashUsecase, mockRatelimitUsecase, mockRecaptchaUsecase)

	url := "abc123"

	mockUrlUsecase.On("GetUrlFromCache", url).Return(nil, nil)
	mockUrlUsecase.On("GetUrlFromDatabase", url).Return(nil, errors.New("url not found"))

	result, err := urlService.ResolveShortUrl(url)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "url not found", err.Error())
	mockUrlUsecase.AssertExpectations(t)
}

func TestAddUrl_Success(t *testing.T) {
	mockUrlUsecase := new(urlU.MockUrlUsecase)
	mockHashUsecase := new(hashU.MockHashUsecase)
	mockRatelimitUsecase := new(ratelimitU.MockRatelimitUsecase)
	mockRecaptchaUsecase := new(recaptchaU.MockRecaptchaUsecase)
	urlService := NewUrlService(mockUrlUsecase, mockHashUsecase, mockRatelimitUsecase, mockRecaptchaUsecase)

	longUrl := "https://example.com"
	ip := "127.0.0.1"
	recaptchaToken := "valid-token"
	shortId := "short123"
	expectedUrl := &entity.Url{LongUrl: longUrl, ShortId: shortId}

	mockRecaptchaUsecase.On("Verify", recaptchaToken).Return(nil)
	mockRatelimitUsecase.On("IsAllowed", ip).Return(true)
	mockHashUsecase.On("GenerateHash").Return(shortId)
	mockUrlUsecase.On("AddUrlToDatabase", mock.MatchedBy(func(url *entity.Url) bool {
		return url.LongUrl == longUrl && url.ShortId == shortId
	})).Return(nil)
	mockRatelimitUsecase.On("Disallow", ip).Return(nil)

	result, err := urlService.AddUrl(longUrl, ip, recaptchaToken)

	assert.NoError(t, err)
	assert.Equal(t, expectedUrl, result)
	mockRecaptchaUsecase.AssertExpectations(t)
	mockRatelimitUsecase.AssertExpectations(t)
	mockHashUsecase.AssertExpectations(t)
	mockUrlUsecase.AssertExpectations(t)
}

func TestAddUrl_InvalidUrlError(t *testing.T) {
	mockUrlUsecase := new(urlU.MockUrlUsecase)
	mockHashUsecase := new(hashU.MockHashUsecase)
	mockRatelimitUsecase := new(ratelimitU.MockRatelimitUsecase)
	mockRecaptchaUsecase := new(recaptchaU.MockRecaptchaUsecase)
	urlService := NewUrlService(mockUrlUsecase, mockHashUsecase, mockRatelimitUsecase, mockRecaptchaUsecase)

	invalidUrl := "invalid-url"
	ip := "127.0.0.1"
	recaptchaToken := "valid-token"

	mockRecaptchaUsecase.On("Verify", recaptchaToken).Return(nil)

	result, err := urlService.AddUrl(invalidUrl, ip, recaptchaToken)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, apperrors.InvalidUrlError.Error())
	mockRecaptchaUsecase.AssertExpectations(t)
}

func TestAddUrl_TooManyRequests(t *testing.T) {
	mockUrlUsecase := new(urlU.MockUrlUsecase)
	mockHashUsecase := new(hashU.MockHashUsecase)
	mockRatelimitUsecase := new(ratelimitU.MockRatelimitUsecase)
	mockRecaptchaUsecase := new(recaptchaU.MockRecaptchaUsecase)
	urlService := NewUrlService(mockUrlUsecase, mockHashUsecase, mockRatelimitUsecase, mockRecaptchaUsecase)

	longUrl := "https://example.com"
	ip := "127.0.0.1"
	recaptchaToken := "valid-token"

	mockRecaptchaUsecase.On("Verify", recaptchaToken).Return(nil)
	mockRatelimitUsecase.On("IsAllowed", ip).Return(false)

	result, err := urlService.AddUrl(longUrl, ip, recaptchaToken)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, apperrors.TooManyRequests.Error())
	mockRecaptchaUsecase.AssertExpectations(t)
	mockRatelimitUsecase.AssertExpectations(t)
}

func TestAddUrl_RecaptchaError(t *testing.T) {
	mockUrlUsecase := new(urlU.MockUrlUsecase)
	mockHashUsecase := new(hashU.MockHashUsecase)
	mockRatelimitUsecase := new(ratelimitU.MockRatelimitUsecase)
	mockRecaptchaUsecase := new(recaptchaU.MockRecaptchaUsecase)
	urlService := NewUrlService(mockUrlUsecase, mockHashUsecase, mockRatelimitUsecase, mockRecaptchaUsecase)

	longUrl := "https://example.com"
	ip := "127.0.0.1"
	recaptchaToken := "invalid-token"

	mockRecaptchaUsecase.On("Verify", recaptchaToken).Return(errors.New("invalid recaptcha"))

	result, err := urlService.AddUrl(longUrl, ip, recaptchaToken)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, apperrors.RecaptchaError.Error())
	mockRecaptchaUsecase.AssertExpectations(t)
}
