// Package main API服务启动入口
package main

import (
	"example/internal/app"
	"example/internal/infrastructure/config"
	"example/internal/infrastructure/depend"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("加载系统配置异常：%+v", err)
	}
	// 注入
	inject, err := depend.NewInjecter(*cfg)
	if err != nil {
		log.Fatalf("[Server]%+v", err)
	}
	if err = app.RunAPIServer(inject); err != nil {
		log.Fatalf("[Server]%+v", err)
	}
}
