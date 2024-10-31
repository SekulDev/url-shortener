package usecase

import "github.com/xinguang/go-recaptcha"

type RecaptchaUsecase interface {
	Verify(token string) error
}

type recaptchaUsecaseImpl struct {
	recaptcha *recaptcha.ReCAPTCHA
}

func NewRecaptchaUsecase(recaptcha *recaptcha.ReCAPTCHA) RecaptchaUsecase {
	return &recaptchaUsecaseImpl{
		recaptcha: recaptcha,
	}
}

func (rl *recaptchaUsecaseImpl) Verify(token string) error {
	return rl.recaptcha.Verify(token)
}
