package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddlewareBuilder struct{}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}
func (*LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		//这两个页面需要验证
		if c.Request.URL.Path == "/users/login" || c.Request.URL.Path == "/users/signup" {
			return
		}
		sess := sessions.Default(c)
		// 验证一下就可以
		if sess.Get("userId") == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
