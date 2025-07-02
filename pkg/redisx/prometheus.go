package redisx

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"net"
	"strconv"
	"time"
)

type PrometheusHook struct {
	vector *prometheus.SummaryVec
}

func NewPrometheusHook(opt prometheus.SummaryOpts) *PrometheusHook {
	vector := prometheus.NewSummaryVec(opt, []string{"cmd", "biz", "key_exists"})
	prometheus.Register(vector)
	return &PrometheusHook{
		vector: vector,
	}
}

func (p *PrometheusHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 相当于，你这里啥也不干
		return next(ctx, network, addr)
	}
}

func (p *PrometheusHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		// 在Redis执行之前
		startTime := time.Now()
		var err error
		defer func() {
			duration := time.Since(startTime).Milliseconds()
			val := ctx.Value("biz")
			biz, ok := val.(string)
			if !ok {
				return
			}
			keyExist := err == redis.Nil
			p.vector.WithLabelValues(
				cmd.Name(),
				biz,
				strconv.FormatBool(keyExist),
			).Observe(float64(duration))
		}()
		// 这个会最终发送命令到 redis 上
		err = next(ctx, cmd)
		// 在 Redis 执行之后
		return err
	}
}

func (p *PrometheusHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}
