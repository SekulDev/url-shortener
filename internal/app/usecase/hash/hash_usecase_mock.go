package usecase

import "github.com/stretchr/testify/mock"

type MockHashUsecase struct {
	mock.Mock
}

func (m *MockHashUsecase) GenerateHash() string {
	args := m.Called()
	return args.String(0)
}
