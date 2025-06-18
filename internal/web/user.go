package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/service"
	"net/http"
)

const emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"

const passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`

type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	svc         service.UserService
	codeSvc     service.CodeService
	jwtHandler
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:         svc,
		codeSvc:     codeSvc,
		jwtHandler:  jwtHandler{},
	}
}

func (h *UserHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/user")
	g.POST("/signup", h.Signup)
	g.POST("/login", h.Login)
	g.POST("/profile", h.Profile)
	g.POST("/logout", h.Logout)
	g.POST("/refresh", h.RefreshToken)
	g.POST("/login_sms", h.LoginSMS)
	g.POST("/login_sms/send", h.SendSMS)
}

func (h *UserHandler) Login(c *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req Req
	if er := c.Bind(&req); er != nil {
		return
	}

	user, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		c.JSON(http.StatusOK, Result{Code: 400, Msg: "账号或者密码错误"})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, Result{Code: 500, Msg: "系统错误"})
		return
	}

	err = h.setJWT(c, user.Id)
	if err != nil {
		c.JSON(http.StatusOK, Result{Code: 500, Msg: "系统错误"})
		return
	}
	c.JSON(http.StatusOK, Result{Code: 200, Msg: "登录成功"})

}

func (h *UserHandler) Signup(c *gin.Context) {
	//接受参数
	type Req struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		RePassword string `json:"rePassword"`
	}
	var req Req
	//Bind方法会根据Context-Type来解析你的数据到req中
	//解析错了，会返回400错误
	if err := c.Bind(&req); err != nil {

		return
	}
	//参数校验
	ok, err := h.emailExp.MatchString(req.Email)
	if err != nil {
		c.JSON(http.StatusOK, Result{Code: 500, Msg: "系统错误"})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, Result{Code: 400, Msg: "邮箱格式不正确"})
		return
	}
	ok, err = h.passwordExp.MatchString(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, Result{Code: 500, Msg: "系统错误"})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, Result{Code: 400, Msg: "密码格式不正确"})
		return
	}
	if req.RePassword != req.Password {
		c.JSON(http.StatusOK, Result{Code: 400, Msg: "两次密码输入不正确"})
		return
	}

	//调用下一层
	err = h.svc.Signup(c.Request.Context(), domain.User{Email: req.Email, Password: req.Password})
	//错误判断
	if err == service.ErrUserDuplicateEmail {
		c.JSON(http.StatusOK, Result{Code: 400, Msg: "邮箱冲突"})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, Result{Code: 500, Msg: "系统错误"})
		return
	}

	c.JSON(http.StatusOK, Result{Code: 200, Msg: "注册成功"})

}

func (h *UserHandler) Profile(c *gin.Context) {
	claims := c.MustGet("claims").(*UserClaims)

	user, err := h.svc.Profile(c.Request.Context(), claims.Uid)

	if err != nil {
		c.JSON(http.StatusOK, Result{Code: 500, Msg: "系统错误"})
		return
	}
	c.JSON(http.StatusOK, Result{Code: 200, Msg: "ok", Data: user})

}

func (h *UserHandler) Logout(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	c.JSON(http.StatusOK, Result{Code: 200, Msg: "退出登录成功"})
}

func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone     string `json:"phone"`
		InputCode string `json:"input_code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := h.codeSvc.Verify(ctx.Request.Context(), "login", req.Phone, req.InputCode)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 500, Msg: "系统错误"})
		return
	}

	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: 400, Msg: "验证码错误"})
		return
	}
	//验证成功
	u, err := h.svc.FindOrCreate(ctx.Request.Context(), req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 500, Msg: "系统错误"})
		return
	}
	//jwt
	err = h.setJWT(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 500, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Code: 200, Msg: "登录成功"})
}

func (h *UserHandler) SendSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	err := h.codeSvc.Send(ctx.Request.Context(), "login", req.Phone)
	if err == service.ErrCodeSendTooMany {
		ctx.JSON(http.StatusOK, Result{Code: 400, Msg: "发送频繁"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 500, Msg: "系统错误"})
	}
	ctx.JSON(http.StatusOK, Result{Code: 200, Msg: "发送成功"})
}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	fmt.Println(123)
	t := ExtractToken(ctx)
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		return RtKey, nil
	})
	fmt.Println(123)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	fmt.Println(234)
	if token == nil || !token.Valid || claims.Uid == 0 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	fmt.Println(456)
	err = h.setJWTToken(ctx, claims.Uid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "刷新成功",
	})
}
