package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
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

		sess := sessions.Default(ctx)

		if sess.Get("userId") == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		now := time.Now().UnixMilli()
		sess.Set("userId", sess.Get("userId"))
		updateTime := sess.Get("updateTime")
		sess.Options(sessions.Options{
			MaxAge: 1800,
		})

		if updateTime == nil {
			sess.Set("updateTime", now)
			sess.Save()
			return
		}

		updateTimeVal, ok := updateTime.(int64)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if now-updateTimeVal > 15*60 {
			ctx.Set("updateTime", now)
			sess.Save()
		}

	}
}
