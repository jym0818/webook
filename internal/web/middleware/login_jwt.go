package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jym/webook/internal/web"
	"net/http"
)

type LoginJWTMiddlewareBuilder struct{}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (*LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/users/login" || c.Request.URL.Path == "/users/signup" ||
			c.Request.URL.Path == "/users/login_sms/code/send" || c.Request.URL.Path == "/users/login_sms" ||
			c.Request.URL.Path == "/oauth2/wechat/authurl" || c.Request.URL.Path == "/oauth2/wechat/callback" ||
			c.Request.URL.Path == "/users/refresh_token" {
			return
		}

		//JWT token
		tokenStr := web.ExtractToken(c)
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

		if claims.UserAgent != c.Request.UserAgent() {

			//严重安全问题,重新登录
			//你是要加监控的
			c.AbortWithStatus(http.StatusUnauthorized)
			return

		}

		c.Set("claims", claims)

	}

}
