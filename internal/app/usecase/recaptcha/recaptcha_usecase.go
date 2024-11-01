package usecase

type RecaptchaUsecase interface {
	Verify(token string) error
}
