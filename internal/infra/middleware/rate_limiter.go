package middleware

import (
	"net/http"
	"strings"

	"github.com/sandronister/mba_challenge1/internal/usecase/limiter"
)

func RateLimiter(l *limiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := getIPRequest(r)
			limit := l.IPLimit
			if token := r.Header.Get("API_KEY"); token != "" {
				key = token
				limit = l.TokenLimit
			}

			allowed, err := l.AllowRequest(r.Context(), key, limit)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if !allowed {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getIPRequest(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}
