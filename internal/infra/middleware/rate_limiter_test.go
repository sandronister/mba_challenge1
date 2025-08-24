package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sandronister/mba_challenge1/internal/infra/internal_error"
	"github.com/sandronister/mba_challenge1/internal/infra/mocks"
	"github.com/sandronister/mba_challenge1/internal/usecase/limiter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateLimiter(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.Repository)
	l := limiter.NewLimiter(mockRepo, 10, 5, 10*time.Second)
	middleware := RateLimiter(l)

	tests := []struct {
		name       string
		apiKey     string
		setup      func()
		wantStatus int
	}{
		{
			name:   "Allow request with IP",
			apiKey: "",
			setup: func() {
				mockRepo.On("AddRequest", ctx, "127.0.0.1", mock.AnythingOfType("time.Duration")).Return(int64(1), (*internal_error.InternalError)(nil)).Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "Block request with IP",
			apiKey: "",
			setup: func() {
				mockRepo.On("AddRequest", ctx, "127.0.0.1", mock.AnythingOfType("time.Duration")).Return(int64(11), (*internal_error.InternalError)(nil)).Once()
			},
			wantStatus: http.StatusTooManyRequests,
		},
		{
			name:   "Allow request with API key",
			apiKey: "valid-api-key",
			setup: func() {
				mockRepo.On("AddRequest", ctx, "valid-api-key", mock.AnythingOfType("time.Duration")).Return(int64(1), (*internal_error.InternalError)(nil)).Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "Block request with API key",
			apiKey: "valid-api-key",
			setup: func() {
				mockRepo.On("AddRequest", ctx, "valid-api-key", mock.AnythingOfType("time.Duration")).Return(int64(6), (*internal_error.InternalError)(nil)).Once()
			},
			wantStatus: http.StatusTooManyRequests,
		},
		{
			name:   "Internal error from repository",
			apiKey: "",
			setup: func() {
				mockRepo.On("AddRequest", ctx, "127.0.0.1", mock.AnythingOfType("time.Duration")).Return(int64(0), &internal_error.InternalError{}).Once()
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			req := httptest.NewRequest("GET", "http://localhost", nil)
			req.RemoteAddr = "127.0.0.1:12345"
			if tt.apiKey != "" {
				req.Header.Set("API_KEY", tt.apiKey)
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			middleware(handler).ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}
