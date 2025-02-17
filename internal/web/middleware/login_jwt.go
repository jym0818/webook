package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jym/webook/internal/web"
	"net/http"
	"strings"
	"time"
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
		claims := &web.UserClaims{}
		t, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("sDKU8mor4FhrCDsFmmMYifqYb8u2X4c7"), nil
		})

		if err != nil {
			//认为没登陆	可能是攻击者伪造的
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//if claims.ExpiresAt.Time.Before(time.Now()) {
		//过期了
		//}
		//err为nil，token为nil
		//不需要校验过期时间，如果过期t.Valid为false
		if t == nil || !t.Valid || claims.Uid == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//登录校验结束----------说明登录成功

		//过了10s 刷新一次
		if claims.ExpiresAt.Time.Sub(time.Now()) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = t.SignedString([]byte("sDKU8mor4FhrCDsFmmMYifqYb8u2X4c7"))
			if err != nil {
				//记录日志 可以不影响程序执行

			}

			c.Header("x-jwt-token", tokenStr)
		}

		c.Set("claims", claims)

	}

}
