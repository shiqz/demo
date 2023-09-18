// Package db 数据库相关
package db

import (
	"demo/internal/infrastructure/config"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"time"
)

// Connector 数据库连接器
type Connector struct {
	Driver *sqlx.DB
	Logger *log.Logger
	debug  bool
}

// Query 查询数据
func (c *Connector) Query(table string) *goqu.SelectDataset {
	return goqu.Dialect("mysql").From(table)
}

// Insert 插入数据
func (c *Connector) Insert(table string) *goqu.InsertDataset {
	return goqu.Dialect("mysql").Insert(table)
}

// Update 更新数据
func (c *Connector) Update(table string) *goqu.UpdateDataset {
	return goqu.Dialect("mysql").Update(table)
}

// PerSQL PerSQL
func (c *Connector) PerSQL(sql string, args []interface{}, err error) (string, []interface{}, error) {
	if c.debug {
		c.Logger.Debugf("[SQL]%s", sql)
	}
	return sql, args, err
}

// NewMySQL 连接MySQL实例
func NewMySQL(cfg config.MySQL, lg *log.Logger) (*Connector, error) {
	dsn := cfg.Builder().FormatDSN()
	lg.Tracef("[db]%s", dsn)
	dc, err := sqlx.Connect("mysql", dsn)
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
	conn := &Connector{Driver: dc, Logger: lg, debug: true}
	return conn, nil
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
