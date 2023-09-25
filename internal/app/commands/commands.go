// Package commands 命令行对应处理器
package commands

import (
	"example/internal/app"
	"example/internal/infrastructure/config"
	"example/internal/infrastructure/depend"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
)

// command 命令行工具
type command struct {
	cfg    config.AppConfig
	inject *depend.Injecter
}

var (
	commander *command // 命令执行者
)

// Init 初始化命令
func Init(cfgFile string) error {
	if cfgFile != "" {
		if err := os.Setenv(config.EnvConfigName, cfgFile); err != nil {
			return errors.Wrap(err, "init config failed")
		}
	}
	cfg, err := config.Init()
	if err != nil {
		return errors.Wrap(err, "init config failed")
	}
	inject, err := depend.NewInjecter(*cfg)
	if err != nil {
		return errors.Wrap(err, "Inject failed")
	}
	commander = &command{*cfg, inject}
	return nil
}

// 启动API服务
func (c *command) runAPIServer() {
	if err := app.RunAPIServer(c.inject); err != nil {
		log.Errorf("[Server]%+v", err)
	}
}

// 查看版本
func (c *command) version() {
	log.Infof("version: %s", app.Version)
}

// 账户控制器
func (c *command) account() *AccountHandler {
	return NewAccountCommand(c.inject)
}

// 账户控制器
func (c *command) user() *UserHandler {
	return NewUserCommand(c.inject)
}
