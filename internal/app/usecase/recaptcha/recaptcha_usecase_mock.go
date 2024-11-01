package usecase

import "github.com/stretchr/testify/mock"

type MockRecaptchaUsecase struct {
	mock.Mock
}

func (m *MockRecaptchaUsecase) Verify(token string) error {
	args := m.Called(token)
	return args.Error(0)
}
