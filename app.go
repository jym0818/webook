package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jym0818/webook/internal/events"
	"github.com/robfig/cron/v3"
)

type App struct {
	web       *gin.Engine
	consumers []events.Consumer
	cron      *cron.Cron
}
