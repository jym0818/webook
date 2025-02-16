package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
		//刷新登录

		updateTime := sess.Get("update_time")
		sess.Set("userId", sess.Get("userId"))
		sess.Options(sessions.Options{
			MaxAge: 30,
		})
		now := time.Now().UnixMilli()
		//说明还没有设置，也就是说是登陆后的第一次请求
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()

			return
		}
		//判断是否需要刷新，因为sess.Get获取的是空接口类型
		updateTimeVal, ok := updateTime.(int64)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if now-updateTimeVal > 15*1000 {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
	}
}
