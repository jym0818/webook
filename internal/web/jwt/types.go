package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	ExtractToken(c *gin.Context) string
	SetLoginToken(c *gin.Context, uid int64) error
	SetJWTToken(c *gin.Context, uid int64, ssid string) error
	ClearToken(c *gin.Context) error
	CheckSession(c *gin.Context, ssid string) error
}
type RefreshClaims struct {
	jwt.RegisteredClaims
	//声明你自己要放入token里面的数据
	Uid  int64
	Ssid string
}

type UserClaims struct {
	//嵌入这个结构体实现了 jwt.Claims接口，从而可以传入函数
	jwt.RegisteredClaims
	//声明你自己要放入token里面的数据
	Uid int64
	//自己随便加 但是最好不要加入敏感数据
	UserAgent string
	Ssid      string
}
