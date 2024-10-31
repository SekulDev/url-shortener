package usecase

import "github.com/xinguang/go-recaptcha"

type RecaptchaUsecase interface {
	Verify(token string) error
}

type RecaptchaUsecaseImpl struct {
	recaptcha *recaptcha.ReCAPTCHA
}

func NewRecaptchaUsecase(recaptcha *recaptcha.ReCAPTCHA) *RecaptchaUsecaseImpl {
	return &RecaptchaUsecaseImpl{
		recaptcha: recaptcha,
	}
}

func (rl *RecaptchaUsecaseImpl) Verify(token string) error {
	return rl.recaptcha.Verify(token)
}
