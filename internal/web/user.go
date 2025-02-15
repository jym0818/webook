package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/service"
	"net/http"
)

const emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"

const passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`

// UserHandler 表示与user相关的路由处理
type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc:         svc,
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

func (u *UserHandler) RegisterRouters(s *gin.Engine) {
	ug := s.Group("/users")
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

func (u *UserHandler) Signup(c *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := c.Bind(&req); err != nil {
		//记录日志 而不是把具体的错误返回给前端
		c.JSON(http.StatusOK, "系统错误")
		return
	}
	//参数校验  使用正则匹配

	ok, err := u.emailExp.MatchString(req.Email)
	//err不为空，说明正则表达式写错了，而不是匹配失败
	if err != nil {
		//记录日志，而不是返回具体错误给前端
		c.JSON(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		c.JSON(http.StatusOK, "邮箱格式不正确")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		//记录日志，而不是返回具体错误给前端
		c.JSON(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		c.JSON(http.StatusOK, "密码必须包含数字、特殊字符，并且长度不能小于 8 位")
		return
	}
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusOK, "两次密码不同")
		return
	}

	//调用svc的方法 下一层service层
	err = u.svc.SignUp(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		c.JSON(http.StatusOK, "重复邮箱，请换一个邮箱")
		return
	}
	if err != nil {
		//记录日志，而不是返回具体错误给前端
		c.JSON(http.StatusOK, "系统错误")
		return
	}

	c.JSON(http.StatusOK, "注册成功")
}
func (u *UserHandler) Login(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := c.Bind(&req); err != nil {
		//记录日志 而不是把具体的错误返回给前端
		c.JSON(http.StatusOK, "系统错误")
		return
	}
	err := u.svc.Login(c, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		c.JSON(http.StatusOK, "账号或者密码错误")
		return
	}
	if err != nil {
		//记录日志 而不是把具体的错误返回给前端
		c.JSON(http.StatusOK, "系统错误")
		return
	}

	c.JSON(http.StatusOK, "登录成功")

}
func (u *UserHandler) Edit(c *gin.Context) {

}
func (u *UserHandler) Profile(c *gin.Context) {

}
