// Package commands 命令行对应处理器
package commands

import (
	"context"
	"demo/internal/app"
	"demo/internal/infrastructure/config"
	"demo/internal/pkg/db"
	"demo/internal/pkg/logger"
	log "github.com/sirupsen/logrus"
	"os"
)

// command 命令行工具
type command struct {
	cfg *config.AppConfig
	lg  *logger.Logger
	dc  *db.Connector
	rdb *db.Redis
}

var (
	commander *command                     // 命令执行者
	depends   = []string{"mysql", "redis"} // 默认依赖
)

// Init 初始化命令
func Init(cfgFile string) {
	if cfgFile != "" {
		if err := os.Setenv(config.EnvConfigName, cfgFile); err != nil {
			log.Fatal("init config failed", err)
		}
	}
	log.SetLevel(log.InfoLevel)
	cfg, err := config.Init()
	if err != nil {
		log.Fatal("init config failed", err)
	}
	lg := logger.New(cfg.Server.LogLevel)
	commander = &command{
		cfg: cfg,
		lg:  lg,
	}
}

// depends 批量初始化依赖, 若未指定，则初始化所有依赖
func (c *command) depends(args ...string) *command {
	if len(args) == 0 {
		args = depends
	}
	if len(args) == 1 {
		return c.depend(args[0])
	}
	for _, name := range args {
		c.depend(name)
	}
	return c
}

// 初始化依赖
func (c *command) depend(name string) *command {
	switch name {
	case "mysql":
		// 初始化 database connection
		dc, err := db.NewMySQL(c.cfg.MySQL, c.lg)
		if err != nil {
			log.Fatalf("[Commander]%+v", err)
		}
		c.dc = dc
	case "redis":
		// 初始化 redis connection
		rdb, err := db.NewRedis(c.cfg.Redis)
		if err != nil {
			log.Fatalf("[Commander]%+v", err)
		}
		c.rdb = rdb
	}
	return c
}

// 启动API服务
func (c *command) runAPIServer() {
	app.RunAPIServer(*c.cfg)
}

// Version 查看版本
func (c *command) version() {
	log.Infof("version: %s", app.Version)
}

// AdminCreate 创建管理员
func (c *command) adminCreate(email, pass, role string) error {
	ctx := context.Background()
	return NewAccountCommand(c.dc, c.rdb, c.lg).Create(ctx, email, pass, role)
}
