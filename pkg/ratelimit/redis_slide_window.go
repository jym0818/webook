package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var luaScript string

type RedisSlideWindow struct {
	cmd redis.Cmdable

	window    time.Duration
	threshold int
}

func NewRedisSlideWindow(cmd redis.Cmdable, window time.Duration, threshold int) *RedisSlideWindow {
	return &RedisSlideWindow{
		cmd:       cmd,
		window:    window,
		threshold: threshold,
	}
}

func (r *RedisSlideWindow) Limit(ctx context.Context, key string) (bool, error) {
	return r.cmd.Eval(ctx, luaScript, []string{key}, time.Now().UnixMilli(), r.window.Milliseconds(), r.threshold).Bool()
}
