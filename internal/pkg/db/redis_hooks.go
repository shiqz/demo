package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

const (
	// SlowKey 执行慢命令时长
	SlowKey = time.Second
)

// HealthHook redis 健康监控钩子
type HealthHook struct {
}

// DialHook redis hook for on dial
func (HealthHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

// ProcessHook redis hook for cmd process
func (HealthHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmd)
		// 记录耗时命令
		if cut := time.Since(start); cut >= SlowKey {
			log.Warnf("[RedisHealth] command: %s, cuttime: %s", cmd.String(), cut.String())
		}
		return err
	}
}

// ProcessPipelineHook redis hook for pipeline
func (HealthHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}
