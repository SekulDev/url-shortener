package usecase

const rateLimitKey = "rate_limit_%s"

type RatelimitUsecase interface {
	IsAllowed(ip string) bool
	Disallow(ip string) error
}
