package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jym0818/webook/internal/events"
)

type App struct {
	web       *gin.Engine
	consumers []events.Consumer
}
