// Package app 启动API服务
package app

import (
	"context"
	"example/internal/app/router"
	"example/internal/infrastructure/config"
	"example/internal/pkg/db"
	"example/internal/pkg/logger"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// Version 版本号
const Version = "v0.2.7"

// APIServer API服务
type APIServer struct {
	Logger *logger.Logger
	mux    chi.Router
}

// ServeHTTP 实现 http Handler 接口
func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Start 启动HTTP服务
func (s *APIServer) Start(ctx context.Context, addr string) error {
	srv := &http.Server{Addr: addr, Handler: s}
	ec := make(chan error, 1)
	s.Logger.Infof("[Server]running on %s", addr)

	go func(ec chan<- error) {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			ec <- errors.WithStack(fmt.Errorf("[Server]shutdonwn: %w", err))
		}
	}(ec)

	shutdown := func(timeout time.Duration) error {
		fmt.Println()
		s.Logger.Info("[Server]shutting...")
		shutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := srv.Shutdown(shutCtx); err != nil {
			return errors.WithStack(fmt.Errorf("[Server]shutdonwn: %w", err))
		}
		return nil
	}
	select {
	case err := <-ec:
		return err
	case <-ctx.Done():
		return shutdown(time.Second * 5)
	}
}

// RunAPIServer 运行API服务
func RunAPIServer(cfg config.AppConfig) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	srv := &APIServer{Logger: logger.New(cfg.Server.LogLevel)}
	srv.Logger.Tracef("[config]%s", cfg.GetConfigFile())

	// init database connection
	dc, err := db.NewMySQL(cfg.MySQL, srv.Logger)
	if err != nil {
		srv.Logger.Fatalf("[Server]failed run: %+v", err)
	}

	// init redis connection
	rdb, err := db.NewRedis(cfg.Redis)
	if err != nil {
		srv.Logger.Fatalf("[Server]failed run: %+v", err)
	}

	srv.mux = router.Init(dc, rdb, srv.Logger)
	if err = srv.Start(ctx, cfg.Server.Addr); err != nil {
		srv.Logger.Fatalf("%+v", err)
	}
}
