package ratelimit

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jym0818/webook/pkg/ratelimit"
	"net/http"
)

type Builder struct {
	limiter ratelimit.Limiter
	prefix  string
}

func NewBuilder(limiter ratelimit.Limiter) *Builder {
	return &Builder{limiter: limiter, prefix: "ip-limiter"}
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
	return b.limiter.Limit(ctx, key)
}
