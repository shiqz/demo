// Package config API服务配置
package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	// EnvConfigName 配置文件对应环境变量名称
	EnvConfigName = "CONFIG_FILE"
	// DefaultConfigType 默认配置文件类型
	DefaultConfigType = "yaml"
	// DefaultConfigFile 默认配置文件
	DefaultConfigFile = "./configs/config.yaml"
)

// AppConfig API服务配置结构
type AppConfig struct {
	Server Server
	MySQL  MySQL
	Redis  Redis
}

// Server Server配置
type Server struct {
	Addr     string `yaml:"addr"`
	Version  string `yaml:"version"`
	LogLevel string `yaml:"logLevel"`
}

// MySQL MySQL数据库配置
type MySQL struct {
	Host            string `yaml:"host"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxIdleTime time.Duration
	ConnMaxLifeTime time.Duration
}

// Redis Redis配置
type Redis struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
}

// DSN 获取MySQL数据库连接
func (c MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", c.User, c.Password, c.Host, c.Database)
}

// Init 初始化配置
func Init() (*AppConfig, error) {
	cfgFile, ok := os.LookupEnv(EnvConfigName)
	if !ok || cfgFile == "" {
		cfgFile = DefaultConfigFile
	}
	viper.SetConfigFile(cfgFile)
	viper.SetConfigType(DefaultConfigType)
	log.Tracef("[config]%s", cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.WithStack(err)
	}
	var cfg AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.WithStack(err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(&cfg); err != nil {
			log.Printf("刷新配置文件异常：%+v", errors.WithStack(err))
		}
	})
	return &cfg, nil
}
