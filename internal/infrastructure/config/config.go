// Package config API服务配置
package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/go-sql-driver/mysql"
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
	file   string
	Server Server
	MySQL  MySQL
	Redis  Redis
}

// GetConfigFile 获取配置文件
func (c *AppConfig) GetConfigFile() string {
	return c.file
}

// Server Server配置
type Server struct {
	Addr     string `yaml:"addr"`
	LogLevel string `yaml:"logLevel"`
}

// MySQL MySQL数据库配置
type MySQL struct {
	Host            string `yaml:"host"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	Debug           bool
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

// Builder 构建MySQL配置
func (c MySQL) Builder() *mysql.Config {
	return &mysql.Config{
		User:    c.User,
		Passwd:  c.Password,
		Net:     "tcp",
		Addr:    c.Host,
		DBName:  c.Database,
		Timeout: time.Second,
	}
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
			log.Errorf("刷新配置文件异常：%+v", errors.WithStack(err))
		}
	})
	cfg.file = cfgFile
	return &cfg, nil
}
