package main

import (
	"context"
	"fmt"

	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sandronister/mba_challenge1/configs"
	"github.com/sandronister/mba_challenge1/configs/logger"
	repository "github.com/sandronister/mba_challenge1/internal/infra/database"
	"github.com/sandronister/mba_challenge1/internal/infra/middleware"
	"github.com/sandronister/mba_challenge1/internal/usecase/limiter"
)

func main() {
	if err := run(); err != nil {
		logger.Fatal(err.Error(), err)
	}
}

func run() (err error) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := configs.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config", err)
	}

	srv := &http.Server{
		Addr:         cfg.ServerPort,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      handler(cfg.RedisAddr, int64(cfg.IPLimit), int64(cfg.TokenLimit), cfg.BlockTime),
	}
	srvErr := make(chan error, 1)

	go func() {
		logger.Info(fmt.Sprintf("Starting server on port '%s'...", cfg.ServerPort[1:]))
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		return
	case <-ctx.Done():
		stop()
	}

	err = srv.Shutdown(ctx)
	return
}

func handler(redisAddr string, ipLimit, tokenLimit int64, blockTime time.Duration) http.Handler {
	repo := repository.NewRedisRepository(redisAddr)
	l := limiter.NewLimiter(repo, ipLimit, tokenLimit, blockTime)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Rate limiter!"))
	})

	return middleware.RateLimiter(l)(mux)
}
