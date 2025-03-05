package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

var (
	ATKey = []byte("sDKU8mor4FhrCDsFmmMYifqYb8u2X4c7")
	RTKey = []byte("sDKU8mor4FhrCDsFmmMYifqYb9u2X4c8")
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd: cmd,
	}
}

func (h *RedisJWTHandler) ExtractToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")

	sges := strings.SplitN(token, " ", 2)

	if len(sges) != 2 {
		return ""
	}
	return sges[1]
}

func (h *RedisJWTHandler) SetLoginToken(c *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.SetJWTToken(c, uid, ssid)
	if err != nil {
		return err
	}
	err = h.setRefreshToken(c, uid, ssid)
	return err
}

func (h *RedisJWTHandler) setRefreshToken(c *gin.Context, uid int64, ssid string) error {
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid:  uid,
		Ssid: ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(RTKey)
	if err != nil {
		c.JSON(http.StatusOK, "系统错误")
		return err
	}
	c.Header("x-refresh-token", tokenStr)
	return nil
}

func (h *RedisJWTHandler) ClearToken(c *gin.Context) error {
	c.Header("x-jwt-token", "")
	c.Header("x-refresh-token", "")
	claims, ok := c.Get("claims")
	if !ok {

		return errors.New("系统错误")
	}
	claimsValue, ok := claims.(*UserClaims)
	if !ok {

		return errors.New("系统错误")
	}
	//过期时间与长token相同
	err := h.cmd.Set(c, fmt.Sprintf("users:ssid:%s", claimsValue.Ssid), "", time.Hour*24*7).Err()
	return err
}

func (h *RedisJWTHandler) CheckSession(c *gin.Context, ssid string) error {
	_, err := h.cmd.Exists(c, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	return err

}

func (h *RedisJWTHandler) SetJWTToken(c *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid:       uid,
		UserAgent: c.Request.UserAgent(),
		Ssid:      ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(ATKey)
	if err != nil {
		c.JSON(http.StatusOK, "系统错误")
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}
