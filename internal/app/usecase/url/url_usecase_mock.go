package usecase

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"url-shortener/internal/domain/entity"
)

type MockUrlUsecase struct {
	mock.Mock
}

func (m *MockUrlUsecase) GetUrlFromCache(url string) (*entity.Url, error) {
	args := m.Called(url)
	result, ok := args.Get(0).(*entity.Url)
	if !ok && result != nil {
		return nil, errors.New("invalid type")
	}
	return result, args.Error(1)
}

func (m *MockUrlUsecase) GetUrlFromDatabase(url string) (*entity.Url, error) {
	args := m.Called(url)
	result, ok := args.Get(0).(*entity.Url)
	if !ok && result != nil {
		return nil, errors.New("invalid type")
	}
	return result, args.Error(1)
}

func (m *MockUrlUsecase) AddUrlToCache(url *entity.Url) error {
	args := m.Called(url)
	return args.Error(0)
}

func (m *MockUrlUsecase) AddUrlToDatabase(url *entity.Url) error {
	args := m.Called(url)
	return args.Error(0)
}
