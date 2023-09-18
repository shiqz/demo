// Package app 启动API服务
package app

import (
	"context"
	"demo/internal/app/router"
	"demo/internal/pkg/db"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// APIServer API服务
type APIServer struct {
	Logger *logrus.Logger
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
			ec <- fmt.Errorf("[Server]shutdonwn: %w", err)
		}
	}(ec)

	shutdown := func(timeout time.Duration) error {
		fmt.Println()
		s.Logger.Info("[Server]shutting...")
		shutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := srv.Shutdown(shutCtx); err != nil {
			return fmt.Errorf("[Server]shutdonwn: %w", err)
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

// NewAPIServer 创建API服务实例
func NewAPIServer(db *db.Connector, rdb *redis.Client, lg *logrus.Logger) *APIServer {
	s := &APIServer{Logger: lg}
	s.mux = router.Init(db, rdb, lg)
	return s
}
