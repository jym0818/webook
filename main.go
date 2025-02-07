package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jym/webook/internal/web"
)

func main() {
	s := gin.Default()
	u := &web.UserHandler{}
	u.RegisterRouters(s)
	s.Run(":8080")
}
