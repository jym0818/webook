//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jym/webook/internal/repository"
	"github.com/jym/webook/internal/repository/cache"
	"github.com/jym/webook/internal/repository/dao"
	"github.com/jym/webook/internal/service"
	"github.com/jym/webook/internal/web"
	"github.com/jym/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,
		dao.NewUserDAO, cache.NewCodeCache, cache.NewUserCache,
		repository.NewUserReposity, repository.NewCodeRepository,
		service.NewUserService, service.NewCodeService,
		ioc.InitSMSService, ioc.InitOAuth2WechatService,
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		//组装gin.Default 和中间件
		ioc.InitGin,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}

//实际上我们使用的wire生成代码中的InitWebServer方法，wire.go会被忽略
