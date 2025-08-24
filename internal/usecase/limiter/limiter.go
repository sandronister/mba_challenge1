package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/sandronister/mba_challenge1/configs/logger"
	repository "github.com/sandronister/mba_challenge1/internal/infra/database"
)

type Limiter struct {
	Repo       repository.Repository
	IPLimit    int64
	TokenLimit int64
	BlockTime  time.Duration
}

func NewLimiter(repository repository.Repository, ipLimit, tokenLimit int64, blockTime time.Duration) *Limiter {
	return &Limiter{
		Repo:       repository,
		IPLimit:    ipLimit,
		TokenLimit: tokenLimit,
		BlockTime:  blockTime,
	}
}

func (l *Limiter) AllowRequest(ctx context.Context, identifier string, limit int64) (bool, error) {
	currentCount, err := l.Repo.AddRequest(ctx, identifier, time.Second)
	if err != nil {
		return false, err
	}

	logger.Info(fmt.Sprintf("Identifier [%30s] -> request %3d/%3d", identifier, currentCount, limit))

	if currentCount > limit {
		logger.Warn(fmt.Sprintf("Identifier [%30s] blocked!", identifier))
		return false, nil
	}

	return true, nil
}
