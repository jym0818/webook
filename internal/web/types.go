package web

import "github.com/gin-gonic/gin"

type Result struct {
	//业务错误码
	//501001 =>代表验证码 5是系统错误 01是用户系统 001是错误的一种
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type handler interface {
	RegisterRouters(s *gin.Engine)
}
