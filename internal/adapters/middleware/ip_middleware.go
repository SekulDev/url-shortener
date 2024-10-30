package middleware

import (
	"context"
	"net/http"
)

type IpMiddleware struct {
}

func NewIpMiddleware() *IpMiddleware {
	return &IpMiddleware{}
}

const IpContextKey string = "ip"

func (m *IpMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			ip = xff
		}

		ctx := context.WithValue(r.Context(), IpContextKey, ip)

		next(w, r.WithContext(ctx))
	}
}
