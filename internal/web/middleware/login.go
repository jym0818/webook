package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jym0818/webook/internal/web"
	"net/http"
	"strings"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		h := ctx.GetHeader("Authorization")
		if h == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(h, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		t := segs[1]
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("sDKU8mor4FhrCDsFmmMYifqYb8u2X4c7"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !token.Valid || claims.Uid == 0 || token == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.UserAgent != ctx.Request.UserAgent() {
			//记录日志
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//刷新登录
		if claims.ExpiresAt.Time.Sub(time.Now()) < time.Minute*15 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
			t, err = token.SignedString([]byte("sDKU8mor4FhrCDsFmmMYifqYb8u2X4c7"))
			if err != nil {
				//记录日志
			}
			ctx.Header("x-jwt-token", t)
		}

		ctx.Set("claims", claims)

	}
}
