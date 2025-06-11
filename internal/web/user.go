package web

import (
	regexp "github.com/dlclark/regexp2"
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
	s.POST("/signup", h.Signup)
}

func (h *UserHandler) Signup(c *gin.Context) {
	//接受参数
	type Req struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		RePassword string `json:"rePassword"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, Result{Code: 500, Msg: "系统错误"})
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
