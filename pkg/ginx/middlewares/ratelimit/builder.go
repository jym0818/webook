package ratelimit

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

//go:embed slide_window.lua
var luaScript string

type Builder struct {
	cmd redis.Cmdable

	window    time.Duration
	threshold int
	prefix    string
}

func NewBuilder(cmd redis.Cmdable, window time.Duration, threshold int) *Builder {
	return &Builder{cmd: cmd, window: window, threshold: threshold, prefix: "ip-limiter"}
}

func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		limited, err := b.limit(c)
		if err != nil {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		if limited {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}

func (b *Builder) limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP())
	return b.cmd.Eval(ctx, luaScript, []string{key}, time.Now().UnixMilli(), b.window.Milliseconds(), b.threshold).Bool()
}
