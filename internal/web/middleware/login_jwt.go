package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type LoginJWTMiddlewareBuilder struct{}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (*LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/users/login" || c.Request.URL.Path == "/users/signup" {
			return
		}

		//JWT token
		token := c.GetHeader("Authorization")
		if token == "" {
			//没登陆
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		sges := strings.SplitN(token, " ", 2)
		//传的格式错误，瞎几把传的，相当于没登陆
		if len(sges) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := sges[1]
		t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("sDKU8mor4FhrCDsFmmMYifqYb8u2X4c7"), nil
		})

		if err != nil {
			//认为没登陆	可能是攻击者伪造的
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//err为nil，token为nil
		if t == nil || !t.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//登录校验结束----------说明登录成功

	}

}
