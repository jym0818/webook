package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type jwtHandler struct {
	//access token key
	atKey []byte
	//refresh token key
	rtKey []byte
}

func newJWTHandler() jwtHandler {
	return jwtHandler{
		atKey: []byte("sDKU8mor4FhrCDsFmmMYifqYb8u2X4c7"),
		rtKey: []byte("sDKU8mor4FhrCDsFmmMYifqYb9u2X4c8")}
}

func (h jwtHandler) setLoginToken(c *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.setJWTToken(c, uid, ssid)
	if err != nil {
		return err
	}
	err = h.setRefreshToken(c, uid, ssid)
	return err
}

func (h jwtHandler) setJWTToken(c *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid:       uid,
		UserAgent: c.Request.UserAgent(),
		Ssid:      ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(h.atKey)
	if err != nil {
		c.JSON(http.StatusOK, "系统错误")
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}

func (h jwtHandler) setRefreshToken(c *gin.Context, uid int64, ssid string) error {
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid:  uid,
		Ssid: ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(h.rtKey)
	if err != nil {
		c.JSON(http.StatusOK, "系统错误")
		return err
	}
	c.Header("x-refresh-token", tokenStr)
	return nil
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

func ExtractToken(c *gin.Context) string {
	//JWT token
	token := c.GetHeader("Authorization")

	sges := strings.SplitN(token, " ", 2)

	if len(sges) != 2 {
		return ""
	}
	return sges[1]
}
