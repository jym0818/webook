package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jym0818/webook/internal/web"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type LoginMiddlewareBuilder struct {
	paths []string
	cmd   redis.Cmdable
}

func NewLoginMiddlewareBuilder(cmd redis.Cmdable) *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{
		cmd: cmd,
	}
}

func (l *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		t := web.ExtractToken(ctx)
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
			return web.AtKey, nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.UserAgent != ctx.Request.UserAgent() {
			//记录日志
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		logout, err := l.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid)).Result()
		if logout > 0 || err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("claims", claims)

	}
}
