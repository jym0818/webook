package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/service"
	"net/http"
	"time"
)

const emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"

const passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
const biz = "login"

// UserHandler 表示与user相关的路由处理
type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		svc:         svc,
		codeSvc:     codeSvc,
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

func (u *UserHandler) RegisterRouters(s *gin.Engine) {
	ug := s.Group("/users")
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.JWTProfile)
	//发送验证码
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
}

func (u *UserHandler) LoginSMS(c *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}

	var req Req

	if err := c.Bind(&req); err != nil {
		c.JSON(200, Result{
			Code: 501001,
			Msg:  "系统错误",
		})
		return
	}
	if req.Phone == "" || req.Code == "" {
		//正常使用正则表达式验证，此处简写
		c.JSON(200, Result{
			Code: 501001,
			Msg:  "输入有误",
		})
		return
	}

	ok, err := u.codeSvc.Verify(c, biz, req.Phone, req.Code)
	if err != nil {
		c.JSON(200, Result{
			Code: 501001,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		c.JSON(200, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}
	//设置登录JWTtoken,可是参数怎么获取呢 uid
	//所以我们创建一个接口 通过手机号查找，如果手机号不存在，我们则要创建新用户
	user, err := u.svc.FindOrCreate(c, req.Phone)
	if err != nil {
		c.JSON(200, Result{
			Code: 501001,
			Msg:  "系统错误",
		})
		return
	}
	if err := u.SetJWTToken(c, user.Id); err != nil {
		c.JSON(200, Result{
			Code: 501001,
			Msg:  "系统错误",
		})
	}
	c.JSON(200, Result{Code: 4, Msg: "验证码校验通过"})
}
func (u *UserHandler) SendLoginSMSCode(c *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}

	var req Req

	if err := c.Bind(&req); err != nil {
		c.JSON(200, "系统错误")
		return
	}
	if req.Phone == "" {
		//正常使用正则表达式验证，此处简写
		c.JSON(200, Result{
			Code: 501001,
			Msg:  "输入有误",
		})
		return
	}

	err := u.codeSvc.Send(c, biz, req.Phone)
	switch err {
	case nil:
		c.JSON(200, Result{Msg: "发送成功"})
	case service.ErrSetCodeTooMany:
		c.JSON(200, Result{Code: 4, Msg: "发送频繁"})
	default:
		c.JSON(200, Result{Code: 5, Msg: "系统错误"})
	}
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
	if err == service.ErrUserDuplicate {
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

func (u *UserHandler) LoginJWT(c *gin.Context) {
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
	user, err := u.svc.Login(c, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		c.JSON(http.StatusOK, "账号或者密码错误")
		return
	}
	if err != nil {
		//记录日志 而不是把具体的错误返回给前端
		c.JSON(http.StatusOK, "系统错误")
		return
	}
	//登录成功，保持登录逻辑
	//使用JWT

	//JWT加入个人数据 比如：userID
	if err := u.SetJWTToken(c, user.Id); err != nil {
		c.JSON(http.StatusOK, "系统错误")
		return
	}
	c.JSON(http.StatusOK, "登录成功")

}

func (u *UserHandler) SetJWTToken(c *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid:       uid,
		UserAgent: c.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("sDKU8mor4FhrCDsFmmMYifqYb8u2X4c7"))
	if err != nil {
		c.JSON(http.StatusOK, "系统错误")
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
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
	user, err := u.svc.Login(c, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		c.JSON(http.StatusOK, "账号或者密码错误")
		return
	}
	if err != nil {
		//记录日志 而不是把具体的错误返回给前端
		c.JSON(http.StatusOK, "系统错误")
		return
	}
	//登录成功，保持登录逻辑
	//从c从获取值
	sess := sessions.Default(c)
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		MaxAge: 30,
	})
	err = sess.Save()
	if err != nil {
		c.JSON(http.StatusOK, "系统错误")
		return
	}
	c.JSON(http.StatusOK, "登录成功")

}

func (u *UserHandler) Logout(c *gin.Context) {
	//登录成功，保持登录逻辑
	//从c从获取值
	sess := sessions.Default(c)
	sess.Get("userId")
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	err := sess.Save()
	if err != nil {
		c.JSON(http.StatusOK, "系统错误")
		return
	}
	c.JSON(http.StatusOK, "退出成功")

}
func (u *UserHandler) Edit(c *gin.Context) {

}
func (u *UserHandler) Profile(c *gin.Context) {

}
func (u *UserHandler) JWTProfile(c *gin.Context) {
	claims, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusOK, "系统错误")
		return
	}
	claimsValue, ok := claims.(*UserClaims)
	if !ok {
		c.JSON(http.StatusOK, "系统错误")
		return
	}
	c.JSON(http.StatusOK, claimsValue.Uid)

}

type UserClaims struct {
	//嵌入这个结构体实现了 jwt.Claims接口，从而可以传入函数
	jwt.RegisteredClaims
	//声明你自己要放入token里面的数据
	Uid int64
	//自己随便加 但是最好不要加入敏感数据
	UserAgent string
}
