package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jym0818/webook/internal/web"
	"github.com/jym0818/webook/internal/web/middleware"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWeb(userHandler *web.UserHandler, mdls []gin.HandlerFunc, wechat *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHandler.RegisterRoutes(server)
	wechat.RegisterRoutes(server)
	return server
}

func InitMiddlware(cmd redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		////限流
		//ratelimit.NewBuilder(cmd, time.Second, 100).Build(),
		middleware.NewLoginMiddlewareBuilder().
			IgnorePath("/user/login").
			IgnorePath("/user/signup").
			IgnorePath("/user/login_sms").
			IgnorePath("/user/login_sms/send").
			IgnorePath("/user/refresh").
			Build(),
	}
}
func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"*"},
		//AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 你不加这个，前端是拿不到的
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		// 是否允许你带 cookie 之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 你的开发环境
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
