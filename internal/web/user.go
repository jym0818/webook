package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:         svc,
	}
}

func (h *UserHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/user")
	g.POST("/signup", h.Signup)
	g.POST("/login", h.Login)
	g.POST("/profile", h.Profile)
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

	sess := sessions.Default(c)
	sess.Set("userId", user.Id)
	err = sess.Save()
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
	sess := sessions.Default(c)
	userId := sess.Get("userId")
	id, ok := userId.(int64)
	if !ok {
		c.JSON(http.StatusOK, Result{Code: 500, Msg: "系统错误"})
		return
	}
	c.JSON(http.StatusOK, Result{Code: 200, Data: id})

}
