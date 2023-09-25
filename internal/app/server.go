// Package app 启动API服务
package app

import (
	"context"
	"example/internal/app/response"
	"example/internal/app/router"
	"example/internal/infrastructure/depend"
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
	inject *depend.Injecter
	mux    chi.Router
}

// ServeHTTP 实现 http Handler 接口
func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w = response.NewResponser(w)
	s.mux.ServeHTTP(w, r)
	_, err := w.(*response.Responser).Render()
	if err != nil {
		s.inject.Log.Errorf("response render, %+v", err)
	}
}

// Start 启动HTTP服务
func (s *APIServer) Start(ctx context.Context, addr string) error {
	srv := &http.Server{Addr: addr, Handler: s}
	ec := make(chan error, 1)
	s.inject.Log.Infof("[Server]running on %s", addr)

	go func(ec chan<- error) {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			ec <- errors.WithStack(fmt.Errorf("[Server]shutdonwn: %w", err))
		}
	}(ec)

	// 关闭服务
	shutdown := func(timeout time.Duration) error {
		fmt.Println()
		s.inject.Log.Info("[Server]shutting...")
		shutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		s.inject.Shutdown()
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
func RunAPIServer(inject *depend.Injecter) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go inject.HealthCheck()
	srv := &APIServer{inject: inject}

	srv.mux = router.Init(inject)
	return srv.Start(ctx, inject.Config().Server.Addr)
}
