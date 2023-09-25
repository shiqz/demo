// Package db 数据库相关
package db

import (
	"context"
	"example/internal/infrastructure/config"
	"example/internal/pkg/logger"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

// Connector 数据库连接器
type Connector struct {
	*sqlx.DB
	lg    *logger.Logger
	debug bool
}

// Debug 开启调试
func (c *Connector) Debug(status bool) *Connector {
	c.debug = status
	return c
}

// Tracef 跟踪调试语句
func (c *Connector) Tracef(sql string, args ...any) {
	if c.debug {
		c.lg.Tracef("[SQL]%s %v", sql, args)
	}
}

// Redis Redis连接器
type Redis struct {
	*redis.Client
}

// NewMySQL 连接MySQL实例
func NewMySQL(cfg config.MySQL, lg *logger.Logger) (*Connector, error) {
	dsn := cfg.Builder().FormatDSN()
	lg.Tracef("[db]%s", dsn)
	dc, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("connection mysql %w", err))
	}
	if err = migrate(dc); err != nil {
		return nil, errors.WithStack(fmt.Errorf("mysql migrate %w", err))
	}
	// 配置连接池
	dc.SetMaxIdleConns(cfg.MaxIdleConn)
	dc.SetMaxOpenConns(cfg.MaxOpenConn)
	dc.SetConnMaxLifetime(cfg.ConnMaxLifeTime)
	dc.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	return &Connector{dc, lg, cfg.Debug}, nil
}

// NewRedis 连接Redis实例
func NewRedis(cfg config.Redis, lg *logger.Logger) (*Redis, error) {
	lg.Tracef("[redis]%s", cfg.Host)
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Host,
		Password:     cfg.Password,
		DB:           cfg.Database,
		WriteTimeout: time.Second,
		ReadTimeout:  time.Second,
		DialTimeout:  time.Second,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, errors.WithStack(fmt.Errorf("connect redis %w", err))
	}
	client.AddHook(HealthHook{})
	return &Redis{client}, nil
}
