package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var lauSlideWindow string

//Redis上的滑动窗口算法限流器实现

type RedisSlidingWindowLimiter struct {
	cmd redis.Cmdable
	//窗口大小
	interval time.Duration
	// 阈值
	rate int
	//interval内允许rate请求
	//1s内允许3000个请求
}

func NewRedisSlidingWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) Limiter {
	return &RedisSlidingWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (r RedisSlidingWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {

	return r.cmd.Eval(ctx, lauSlideWindow, []string{key},
		r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()

}
