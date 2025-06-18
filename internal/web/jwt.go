package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

var (
	AtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	RtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvfx")
)

type jwtHandler struct{}

func (h jwtHandler) setJWT(c *gin.Context, uid int64) error {
	err := h.setJWTToken(c, uid)
	if err != nil {
		return err
	}
	return h.setRefreshToken(c, uid)
}

func (h jwtHandler) setJWTToken(c *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       uid,
		UserAgent: c.Request.UserAgent(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := t.SignedString(AtKey)
	if err != nil {
		return err
	}
	c.Header("x-jwt-token", token)
	return nil
}

func (h jwtHandler) setRefreshToken(c *gin.Context, uid int64) error {
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid: uid,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := t.SignedString(RtKey)
	if err != nil {
		return err
	}
	c.Header("x-refresh-token", token)
	return nil
}

func ExtractToken(ctx *gin.Context) string {
	t := ctx.GetHeader("Authorization")

	segs := strings.Split(t, " ")
	if len(segs) != 2 {
		return ""
	}
	tokenStr := segs[1]
	return tokenStr
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid int64
}
