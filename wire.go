//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jym0818/webook/internal/repository"
	"github.com/jym0818/webook/internal/repository/cache"
	"github.com/jym0818/webook/internal/repository/dao"
	"github.com/jym0818/webook/internal/service"
	"github.com/jym0818/webook/internal/web"
	"github.com/jym0818/webook/ioc"
)

var UserService = wire.NewSet(
	cache.NewuserCache,
	dao.NewuserDAO,
	repository.NewuserRepository,
	service.NewuserService,
)

var CodeService = wire.NewSet(
	cache.NewcodeCache,
	repository.NewcodeRepository,
	service.NewcodeService,
)

func InitServer() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitSMS,

		UserService,
		CodeService,
		web.NewUserHandler,
		ioc.InitWeb,
		ioc.InitMiddlware,
	)
	return new(gin.Engine)
}
