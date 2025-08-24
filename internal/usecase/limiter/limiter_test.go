package limiter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sandronister/mba_challenge1/internal/infra/internal_error"
	"github.com/sandronister/mba_challenge1/internal/infra/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLimiter_AllowRequest(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.Repository)
	limiter := NewLimiter(mockRepo, 10, 5, 10*time.Second)

	tests := []struct {
		name       string
		identifier string
		limit      int64
		setup      func()
		want       bool
		wantErr    bool
	}{
		{
			name:       "New Request",
			identifier: "127.0.0.1",
			limit:      10,
			setup: func() {
				mockRepo.On("AddRequest", ctx, "127.0.0.1", mock.AnythingOfType("time.Duration")).Return(int64(1), (*internal_error.InternalError)(nil)).Once()
			},
			want:    true,
			wantErr: false,
		},
		{
			name:       "Within Limit",
			identifier: "127.0.0.1",
			limit:      10,
			setup: func() {
				mockRepo.On("AddRequest", ctx, "127.0.0.1", mock.AnythingOfType("time.Duration")).Return(int64(5), (*internal_error.InternalError)(nil)).Once()
			},
			want:    true,
			wantErr: false,
		},
		{
			name:       "Over Limit",
			identifier: "127.0.0.1",
			limit:      10,
			setup: func() {
				mockRepo.On("AddRequest", ctx, "127.0.0.1", mock.AnythingOfType("time.Duration")).Return(int64(11), (*internal_error.InternalError)(nil)).Once()
			},
			want:    false,
			wantErr: false,
		},
		{
			name:       "Repository Error",
			identifier: "127.0.0.1",
			limit:      10,
			setup: func() {
				mockRepo.On("AddRequest", ctx, "127.0.0.1", mock.AnythingOfType("time.Duration")).Return(int64(-1), &internal_error.InternalError{}).Once()
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got, err := limiter.AllowRequest(ctx, tt.identifier, tt.limit)
			if tt.wantErr {
				assert.Error(t, errors.New(err.Error()))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
			mockRepo.AssertExpectations(t)
		})
	}
}
