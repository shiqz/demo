// Package main API服务启动入口
package main

import (
	"demo/internal/app"
	"demo/internal/infrastructure/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("加载系统配置异常：%+v", err)
	}
	app.RunAPIServer(*cfg)
}
