package usecase

import "github.com/stretchr/testify/mock"

type MockRatelimitUsecase struct {
	mock.Mock
}

func (m *MockRatelimitUsecase) IsAllowed(ip string) bool {
	args := m.Called(ip)
	return args.Bool(0)
}

func (m *MockRatelimitUsecase) Disallow(ip string) error {
	args := m.Called(ip)
	return args.Error(0)
}
