package repository

import (
	"errors"
	"url-shortener/internal/domain/entity"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUrlRepository struct {
	mock.Mock
}

func (m *MockUrlRepository) Create(url *entity.Url) (primitive.ObjectID, error) {
	args := m.Called(url)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *MockUrlRepository) GetByID(id string) (*entity.Url, error) {
	args := m.Called(id)
	url, ok := args.Get(0).(*entity.Url)
	if !ok {
		return nil, errors.New("invalid type")
	}
	return url, args.Error(1)
}

func (m *MockUrlRepository) GetByShortUrl(url string) (*entity.Url, error) {
	args := m.Called(url)
	result, ok := args.Get(0).(*entity.Url)
	if !ok {
		return nil, errors.New("invalid type")
	}
	return result, args.Error(1)
}
