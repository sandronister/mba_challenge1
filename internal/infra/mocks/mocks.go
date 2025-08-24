package mocks

import (
	"context"
	"time"

	"github.com/sandronister/mba_challenge1/internal/infra/internal_error"
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (m *Repository) AddRequest(ctx context.Context, ip string, expiration time.Duration) (int64, *internal_error.InternalError) {
	args := m.Called(ctx, ip, expiration)
	return args.Get(0).(int64), args.Get(1).(*internal_error.InternalError)
}

func (m *Repository) Block(ctx context.Context, ip string, expiration time.Duration) *internal_error.InternalError {
	args := m.Called(ctx, ip, expiration)
	return args.Get(0).(*internal_error.InternalError)
}

func (m *Repository) IsBlocked(ctx context.Context, ip string) (bool, *internal_error.InternalError) {
	args := m.Called(ctx, ip)
	return args.Get(0).(bool), args.Get(1).(*internal_error.InternalError)
}
