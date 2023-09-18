// Package main API服务启动入口
package main

import (
	"context"
	"demo/internal/app"
	"demo/internal/infrastructure/config"
	"demo/internal/pkg/db"
	"demo/internal/pkg/logger"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("加载系统配置异常：%+v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err = run(ctx, *cfg); err != nil {
		log.Fatalf("[Server]failed running: %v", err)
	}
	log.Info("[Server]closed")
}

func run(ctx context.Context, cfg config.AppConfig) error {
	lg := logger.New(cfg.Server.LogLevel)
	// 初始化 database connection
	dc, err := db.NewMySQL(cfg.MySQL, lg)
	if err != nil {
		return err
	}
	// 初始化 redis connection
	rdb := db.NewRedis(cfg.Redis)
	if err = rdb.Ping(ctx).Err(); err != nil {
		return errors.Wrap(err, "connect redis")
	}
	return app.NewAPIServer(dc, rdb, lg).Start(ctx, cfg.Server.Addr)
}
