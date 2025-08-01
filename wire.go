//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jym0818/webook/internal/events/article"
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

var ArticleService = wire.NewSet(
	dao.NewarticleDAO,
	cache.NewarticleCache,
	repository.NewarticleRepository,
	service.NewarticleService,
)

func InitServer() *App {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitSMS,

		UserService,
		CodeService,
		web.NewUserHandler,
		ioc.InitWeb,
		ioc.InitMiddlware,

		web.NewOAuth2WechatHandler,
		ioc.InitWechat,
		ioc.InitWechatCfg,

		ArticleService,
		web.NewArticleHandler,

		ioc.InitKafka,
		ioc.InitKafkaProducer,
		article.NewKafkaProducer,

		service.NewBatchRankingService,
		ioc.InitRankingJob,
		ioc.InitCronJob,
		repository.NewCachedRankingRepository,
		cache.NewRankingRedisCache,
		cache.NewRankingLocalCache,
		ioc.InitIntrGRPCClient,
		ioc.InitEtcd,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
