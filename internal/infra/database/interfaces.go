package database

import (
	"context"
	"time"

	"github.com/sandronister/mba_challenge1/internal/infra/internal_error"
)

type Repository interface {
	AddRequest(ctx context.Context, identifier string, expiration time.Duration) (int64, *internal_error.InternalError)
	Block(ctx context.Context, identifier string, expiration time.Duration) *internal_error.InternalError
	IsBlocked(ctx context.Context, identifier string) (bool, *internal_error.InternalError)
}
