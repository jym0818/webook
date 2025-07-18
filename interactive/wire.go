//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jym0818/webook/interactive/events"
	"github.com/jym0818/webook/interactive/grpc"
	"github.com/jym0818/webook/interactive/ioc"
	"github.com/jym0818/webook/interactive/repository"
	"github.com/jym0818/webook/interactive/repository/cache"
	"github.com/jym0818/webook/interactive/repository/dao"
	"github.com/jym0818/webook/interactive/service"
)

func InitApp() *App {
	wire.Build(
		ioc.InitRedis,
		ioc.InitDB,
		ioc.InitKafka,
		ioc.NewConsumers,
		events.NewReadEventArticleConsumer,
		repository.NewinteractiveRepository,
		cache.NewinteractiveCache,
		dao.NewinteractiveDAO,
		service.NewinteractiveService,
		ioc.InitGRPCxServer,
		grpc.NewInteractiveServiceServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
