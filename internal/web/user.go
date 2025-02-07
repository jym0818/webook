package web

import "github.com/gin-gonic/gin"

// UserHandler 表示与user相关的路由处理
type UserHandler struct {
}

func (u *UserHandler) RegisterRouters(s *gin.Engine) {
	s.POST("/users/signup", u.Signup)
	s.POST("/users/login", u.Login)
	s.POST("/users/edit", u.Edit)
	s.GET("/users/profile", u.Profile)
}

func (u *UserHandler) Signup(c *gin.Context) {

}
func (u *UserHandler) Login(c *gin.Context) {

}
func (u *UserHandler) Edit(c *gin.Context) {

}
func (u *UserHandler) Profile(c *gin.Context) {

}
