package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type jwtHandler struct {
}

func (h jwtHandler) SetJWTToken(c *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid:       uid,
		UserAgent: c.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("sDKU8mor4FhrCDsFmmMYifqYb8u2X4c7"))
	if err != nil {
		c.JSON(http.StatusOK, "系统错误")
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}

type UserClaims struct {
	//嵌入这个结构体实现了 jwt.Claims接口，从而可以传入函数
	jwt.RegisteredClaims
	//声明你自己要放入token里面的数据
	Uid int64
	//自己随便加 但是最好不要加入敏感数据
	UserAgent string
}
