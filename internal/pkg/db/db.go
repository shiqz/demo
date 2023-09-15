// Package db 数据库相关
package db

import (
	"demo/internal/infrastructure/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"time"
)

// NewMySQL 连接MySQL实例
func NewMySQL(cfg config.MySQL) (*sqlx.DB, error) {
	log.Tracef("[db]%s", cfg.DSN())
	dc, err := sqlx.Connect("mysql", cfg.DSN())
	if err != nil {
		return nil, errors.Wrap(err, "connection mysql")
	}
	if err = migrate(dc); err != nil {
		return nil, err
	}
	// 配置连接池
	dc.SetMaxIdleConns(cfg.MaxIdleConn)
	dc.SetMaxOpenConns(cfg.MaxOpenConn)
	dc.SetConnMaxLifetime(cfg.ConnMaxLifeTime)
	dc.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	return dc, nil
}

// NewRedis 连接Redis实例
func NewRedis(cfg config.Redis) *redis.Client {
	log.Tracef("[redis]%s", cfg.Host)
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Host,
		Password:     cfg.Password,
		DB:           cfg.Database,
		WriteTimeout: time.Second,
		ReadTimeout:  time.Second,
	})
	client.AddHook(HealthHook{})
	return client
}
