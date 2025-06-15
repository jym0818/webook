package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jym0818/webook/internal/repository"
	"github.com/jym0818/webook/internal/repository/cache"
	"github.com/jym0818/webook/internal/repository/dao"
	"github.com/jym0818/webook/internal/service"
	"github.com/jym0818/webook/internal/web"
	"github.com/jym0818/webook/internal/web/middleware"
	"github.com/jym0818/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	server := gin.Default()
	// 使用 CORS 中间件处理跨域问题，配置 CORS 参数
	server.Use(cors.New(cors.Config{
		// 允许的源地址（CORS中的Access-Control-Allow-Origin）
		// AllowOrigins: []string{"https://foo.com"},
		// 允许的 HTTP 方法（CORS中的Access-Control-Allow-Methods）
		//如果省略，那么所有方法都允许
		//AllowMethods: []string{"PUT", "PATCH"},
		// 允许的 HTTP 头部（CORS中的Access-Control-Allow-Headers）
		AllowHeaders: []string{"Origin"},
		// 暴露的 HTTP 头部（CORS中的Access-Control-Expose-Headers）
		ExposeHeaders: []string{"Content-Length", "x-jwt-token"},
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

	cmd := redis.NewClient(&redis.Options{
		Addr: "118.25.44.1:6379",
	})
	//限流
	server.Use(ratelimit.NewBuilder(cmd, time.Second, 100).Build())

	//检验登录状态
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePath("/user/login").
		IgnorePath("/user/signup").
		Build())

	//依赖注入
	db, err := gorm.Open(mysql.Open("root:root@tcp(118.25.44.1:13316)/webook"))
	if err != nil {
		panic(err)
	}
	err = dao.InitDB(db)
	if err != nil {
		panic(err)
	}
	userDAO := dao.NewuserDAO(db)
	userCache := cache.NewuserCache(cmd)
	userRepo := repository.NewuserRepository(userDAO, userCache)
	userSvc := service.NewuserService(userRepo)
	userHandler := web.NewUserHandler(userSvc)
	userHandler.RegisterRoutes(server)
	server.Run(":8080")
}
