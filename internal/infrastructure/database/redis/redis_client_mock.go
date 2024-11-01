package database

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedis) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockRedis) Set(key string, value interface{}, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *MockRedis) Exists(keys ...string) (int64, error) {
	args := m.Called(mock.Anything)
	return args.Get(0).(int64), args.Error(1)
}
