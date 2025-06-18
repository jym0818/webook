package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jym0818/webook/internal/service"
	"github.com/jym0818/webook/internal/service/oauth2/wechat"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"time"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	jwtHandler
	stateKey []byte
	cfg      Config
}
type Config struct {
	Secure bool
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService, cfg Config) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:        svc,
		userSvc:    userSvc,
		jwtHandler: jwtHandler{},
		stateKey:   []byte("12345678912345678912345678912345"),
		cfg:        cfg,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthURL(ctx.Request.Context(), state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 500, Msg: "构造扫码登录URL失败"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{

			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.SetCookie("jwt-state", tokenStr, 60*10,
		"/oauth2/wechat/callback", "",
		h.cfg.Secure, true)

	ctx.JSON(http.StatusOK, Result{Code: 500, Data: url})
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		//正常不会走这里  做好监控
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	var sc StateClaims
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !token.Valid {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "登录失败",
		})
		return
	}

	if sc.State != state {
		//记录日志
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "登录失败",
		})
		return
	}

	info, err := h.svc.VerifyCode(ctx.Request.Context(), code, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	//登录成功了
	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	//jwt提取出来
	err = h.setJWT(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "ok",
	})
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}
