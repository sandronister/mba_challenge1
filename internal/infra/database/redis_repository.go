package database

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sandronister/mba_challenge1/configs/logger"
	"github.com/sandronister/mba_challenge1/internal/infra/internal_error"
)

type RedisRepository struct {
	rdb         *redis.Client
	keyRequests string
	keyBlock    string
}

func NewRedisRepository(addr string) *RedisRepository {
	return &RedisRepository{
		rdb: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
		keyRequests: "request:ip:",
		keyBlock:    "block:ip:",
	}
}

func (r *RedisRepository) AddRequest(ctx context.Context, identifier string, expiration time.Duration) (int64, *internal_error.InternalError) {
	script := redis.NewScript(`
        local current
        current = redis.call("GET", KEYS[1])
        if not current then
            redis.call("SET", KEYS[1], 1, "EX", ARGV[1])
            return 1
        else
            current = tonumber(current) + 1
            redis.call("INCR", KEYS[1])
            return current
        end
    `)

	result, err := script.Run(ctx, r.rdb, []string{r.getKeyRequests(identifier)}, expiration.Seconds()).Int()
	if err != nil {
		logger.Error(err.Error(), err)
		return 0, internal_error.NewInternalServerError(err.Error())
	}
	return int64(result), nil
}

func (r *RedisRepository) Block(ctx context.Context, identifier string, expiration time.Duration) *internal_error.InternalError {
	if err := r.rdb.Set(ctx, r.getKeyBlock(identifier), 1, expiration).Err(); err != nil {
		logger.Error(err.Error(), err)
		return internal_error.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *RedisRepository) IsBlocked(ctx context.Context, identifier string) (bool, *internal_error.InternalError) {
	val, err := r.rdb.Get(ctx, r.getKeyBlock(identifier)).Result()
	if errors.Is(err, redis.Nil) {
		return false, internal_error.NewNotFoundError("ip block not found")
	} else if err != nil {
		logger.Error(err.Error(), err)
		return false, internal_error.NewInternalServerError(err.Error())
	}
	return val == "1", nil
}

func (r *RedisRepository) getKeyBlock(identifier string) string {
	return r.keyBlock + identifier
}

func (r *RedisRepository) getKeyRequests(identifier string) string {
	return r.keyRequests + identifier
}
