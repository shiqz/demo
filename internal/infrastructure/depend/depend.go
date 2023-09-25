// Package depend 依赖注入
package depend

import (
	"context"
	"example/internal/infrastructure/config"
	"example/internal/pkg/db"
	"example/internal/pkg/logger"
	"time"
)

// Injecter 依赖注入器
type Injecter struct {
	cfg   config.AppConfig
	conn  *db.Connector
	redis *db.Redis
	Log   *logger.Logger
}

// NewInjecter 创建注入器
func NewInjecter(cfg config.AppConfig) (*Injecter, error) {
	// init logger
	lg := logger.New(cfg.Server.LogLevel)
	lg.Tracef("[config]%s", cfg.GetConfigFile())

	// init database connection
	conn, err := db.NewMySQL(cfg.MySQL, lg)
	if err != nil {
		return nil, err
	}

	// init redis connection
	rdb, err := db.NewRedis(cfg.Redis, lg)
	if err != nil {
		return nil, err
	}

	return &Injecter{cfg, conn, rdb, lg}, nil
}

// Config 获取配置
func (j *Injecter) Config() config.AppConfig {
	return j.cfg
}

// GetDB 获取数据库
func (j *Injecter) GetDB() *db.Connector {
	if j.conn == nil {
		j.Log.Warn("[db]disconnected")
	}
	return j.conn
}

// GetRedis 获取 Redis
func (j *Injecter) GetRedis() *db.Redis {
	if j.conn == nil {
		j.Log.Warn("[Redis]disconnected")
	}
	return j.redis
}

// Shutdown 关闭
func (j *Injecter) Shutdown() {
	_ = j.conn.Close()
	_ = j.redis.Close()
}

// HealthCheck 健康监控
func (j *Injecter) HealthCheck() {
	for range time.Tick(10 * time.Minute) {
		go j.checkConn()
		go j.checkRedis()
	}
}

// checkConn 检查conn
func (j *Injecter) checkConn() {
	var err error
	retry := j.conn == nil
	if retry || j.conn.Ping() != nil {
		if j.conn, err = db.NewMySQL(j.cfg.MySQL, j.Log); err != nil {
			j.Log.Errorf("[HealthCheck]mysql retry connect %+v", err)
		} else {
			j.Log.Info("[HealthCheck]mysql retry connect ok")
		}
	} else {
		j.Log.Info("[HealthCheck]mysql ok")
	}
}

// checkRedis 检查redis
func (j *Injecter) checkRedis() {
	var err error
	retry := j.redis == nil
	if retry || j.redis.Ping(context.Background()).Err() != nil {
		if j.redis, err = db.NewRedis(j.cfg.Redis, j.Log); err != nil {
			j.Log.Errorf("[HealthCheck]redis retry connect %+v", err)
		} else {
			j.Log.Info("[HealthCheck]redis retry connect ok")
		}
	} else {
		j.Log.Info("[HealthCheck]redis ok")
	}
}
