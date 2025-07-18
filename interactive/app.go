package main

import (
	"github.com/jym0818/webook/pkg/grpcx"
	"github.com/jym0818/webook/pkg/saramax"
)

type App struct {
	server    *grpcx.Server
	consumers []saramax.Consumer
}
