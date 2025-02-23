package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	s := InitWebServer()
	s.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})
	s.Run(":8080")
}
