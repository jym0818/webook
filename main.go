package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/jym/webook/internal/repository"
	"github.com/jym/webook/internal/repository/dao"
	"github.com/jym/webook/internal/service"
	"github.com/jym/webook/internal/web"
	"github.com/jym/webook/internal/web/middleware"
	"github.com/jym/webook/pkg/ginx/middleware/ratelimit"
	rd "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	//初始化db和建表
	db := initDB()
	//初始化web
	s := initWebServer()
	//初始化路由
	u := initUser(db)
	//路由注册
	u.RegisterRouters(s)

	s.Run(":8080")
}

func initWebServer() *gin.Engine {
	s := gin.Default()
	rdClient := rd.NewClient(&rd.Options{
		//对应k8s-redis-service.yaml port
		Addr: "webook-redis:6379",
	})
	//限流middleware
	//限流同一IP 1s 100次
	s.Use(ratelimit.NewBuilder(rdClient, time.Second, 100).Build())

	// 使用 CORS 中间件处理跨域问题，配置 CORS 参数
	s.Use(cors.New(cors.Config{
		// 允许的源地址（CORS中的Access-Control-Allow-Origin）
		// AllowOrigins: []string{"https://foo.com"},
		// 允许的 HTTP 方法（CORS中的Access-Control-Allow-Methods）
		//如果省略，那么所有方法都允许
		AllowMethods: []string{"PUT", "PATCH"},
		// 允许的 HTTP 头部（CORS中的Access-Control-Allow-Headers）
		AllowHeaders: []string{"Content-Type", "Authorization"},
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

	// 1. 创建一个基于memstore存的存储
	store, err := redis.NewStore(16, "tcp", "webook-redis:6379", "",
		[]byte("sDKU8mor4FhrCDsFmmMYifqYb8u2X4c7"), []byte("zP6Va4QtIFrAVbDFOzMQwKtXDfZFpM5i"))
	if err != nil {
		panic(err)
	}

	// 2. 注册会话中间件，所有请求的会话将被命名为 "mysession"，数据存储在 Cookie 中
	s.Use(sessions.Sessions("mysession", store))

	s.Use(middleware.NewLoginJWTMiddlewareBuilder().CheckLogin())

	return s
}

func initUser(db *gorm.DB) *web.UserHandler {
	//也就是说我们都是通过New函数创建结构体，而不是手动创建一个结构体  使用依赖注入

	userDao := dao.NewUserDAO(db)               //需要db
	repo := repository.NewUserReposity(userDao) //需要dao，上一层
	svc := service.NewUserService(repo)         //需要repo，上一层
	u := web.NewUserHandler(svc)                //需要service才能初始化handler，上一层
	return u
}
func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(webook-mysql:3306)/webook?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		//我只会在初始化过程中panic
		//一旦panic goroutine就会结束
		//一旦初始化过程出错，应用就不要启动了，所以panic
		panic(err)
	}
	//初始化建表---------实际工作中不会使用这种方法
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
