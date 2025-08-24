package database

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sandronister/mba_challenge1/internal/infra/internal_error"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_AddRequest(t *testing.T) {
	ctx := context.Background()
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer s.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	repo := &RedisRepository{
		rdb:         rdb,
		keyRequests: "request:ip:",
		keyBlock:    "block:ip:",
	}

	tests := []struct {
		name      string
		ip        string
		blockTime time.Duration
		setup     func()
		wantCount int64
		wantErr   *internal_error.InternalError
	}{
		{
			name:      "First Request",
			ip:        "127.0.0.1",
			blockTime: time.Minute,
			setup:     func() {},
			wantCount: 1,
			wantErr:   nil,
		},
		{
			name:      "Increment Request Count",
			ip:        "127.0.0.1",
			blockTime: time.Minute,
			setup: func() {
				_ = s.Set("request:ip:127.0.0.1", "5")
			},
			wantCount: 6,
			wantErr:   nil,
		},
		{
			name:      "Redis Error",
			ip:        "127.0.0.1",
			blockTime: time.Minute,
			setup: func() {
				_ = rdb.Close()
			},
			wantCount: -1,
			wantErr:   internal_error.NewInternalServerError("redis: client is closed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			count, err := repo.AddRequest(ctx, tt.ip, tt.blockTime)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
				assert.Equal(t, tt.wantCount, count)
			}
		})
	}
}

func TestRedisRepository_Block(t *testing.T) {
	ctx := context.Background()
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer s.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	repo := &RedisRepository{
		rdb:         rdb,
		keyRequests: "request:ip:",
		keyBlock:    "block:ip:",
	}

	t.Run("Block IP", func(t *testing.T) {
		err := repo.Block(ctx, "127.0.0.1", time.Minute)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		blocked := s.Exists("block:ip:127.0.0.1")
		assert.True(t, blocked)
	})
}

func TestRedisRepository_IsBlocked(t *testing.T) {
	ctx := context.Background()
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer s.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	repo := &RedisRepository{
		rdb:         rdb,
		keyRequests: "request:ip:",
		keyBlock:    "block:ip:",
	}

	tests := []struct {
		name    string
		ip      string
		setup   func()
		want    bool
		wantErr *internal_error.InternalError
	}{
		{
			name: "IP Not Blocked",
			ip:   "127.0.0.1",
			setup: func() {
				s.Del("block:ip:127.0.0.1")
			},
			want:    false,
			wantErr: internal_error.NewNotFoundError("ip block not found"),
		},
		{
			name: "IP Blocked",
			ip:   "127.0.0.1",
			setup: func() {
				_ = s.Set("block:ip:127.0.0.1", "1")
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "Redis Error",
			ip:   "127.0.0.1",
			setup: func() {
				_ = rdb.Close()
			},
			want:    false,
			wantErr: internal_error.NewInternalServerError("redis: client is closed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			blocked, err := repo.IsBlocked(ctx, tt.ip)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
				assert.Equal(t, tt.want, blocked)
			}
		})
	}
}
