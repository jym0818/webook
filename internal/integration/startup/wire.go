//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jym/webook/internal/repository"
	"github.com/jym/webook/internal/repository/article"
	"github.com/jym/webook/internal/repository/cache"
	"github.com/jym/webook/internal/repository/dao"
	article2 "github.com/jym/webook/internal/repository/dao/article"
	"github.com/jym/webook/internal/service"
	"github.com/jym/webook/internal/web"
	"github.com/jym/webook/internal/web/jwt"
	"github.com/jym/webook/ioc"
)

var thirdProvider = wire.NewSet(InitRedis, InitDB, InitLogger)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,
		dao.NewUserDAO, cache.NewCodeCache, cache.NewUserCache, article2.NewGORMArticleDAO,
		repository.NewUserReposity, repository.NewCodeRepository, article.NewCachedArticleRepository,
		service.NewUserService, service.NewCodeService, service.NewArticleService,
		ioc.InitSMSService, ioc.InitOAuth2WechatService, ioc.NewWechatHandler,
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		//组装gin.Default 和中间件
		ioc.InitGin,
		ioc.InitMiddlewares,
		ioc.InitLogger,
		jwt.NewRedisJWTHandler,
	)
	return new(gin.Engine)
}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(thirdProvider, service.NewArticleService, web.NewArticleHandler, article.NewCachedArticleRepository, article2.NewGORMArticleDAO)
	return &web.ArticleHandler{}
}
