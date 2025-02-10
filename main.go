package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jym/webook/internal/web"
	"strings"
	"time"
)

func main() {
	s := gin.Default()
	// 使用 CORS 中间件处理跨域问题，配置 CORS 参数
	s.Use(cors.New(cors.Config{
		// 允许的源地址（CORS中的Access-Control-Allow-Origin）
		// AllowOrigins: []string{"https://foo.com"},
		// 允许的 HTTP 方法（CORS中的Access-Control-Allow-Methods）
		//如果省略，那么所有方法都允许
		AllowMethods: []string{"PUT", "PATCH"},
		// 允许的 HTTP 头部（CORS中的Access-Control-Allow-Headers）
		AllowHeaders: []string{"Content-Type"},
		// 暴露的 HTTP 头部（CORS中的Access-Control-Expose-Headers）
		ExposeHeaders: []string{"Content-Length"},
		// 是否允许携带身份凭证（CORS中的Access-Control-Allow-Credentials）
		AllowCredentials: true,
		// 允许源的自定义判断函数，返回true表示允许，false表示不允许
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 允许你的开发环境
				return true
			}
			// 允许包含 "yourcompany.com" 的源
			return strings.Contains(origin, "yourcompany.com")
		},
		// 用于缓存预检请求结果的最大时间（CORS中的Access-Control-Max-Age）
		MaxAge: 12 * time.Hour,
	}))
	u := web.NewUserHandler()
	u.RegisterRouters(s)
	s.Run(":8080")
}
