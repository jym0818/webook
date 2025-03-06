package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/service"
	ijwt "github.com/jym/webook/internal/web/jwt"
	"github.com/jym/webook/pkg/logger"
	"net/http"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.LoggerV1
}

func NewArticleHandler(svc service.ArticleService, l logger.LoggerV1) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}

func (h *ArticleHandler) RegisterRouters(s *gin.Engine) {
	g := s.Group("/articles")
	g.POST("/edit", h.Edit)
}

func (h *ArticleHandler) Edit(c *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}
	//检测输入
	user, _ := c.Get("claims")
	claims, ok := user.(*ijwt.UserClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		h.l.Error("未发现用户的session信息")
		return
	}

	//调用service
	id, err := h.svc.Save(c, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
	if err != nil {
		c.JSON(200, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("保存帖子失败", logger.Error(err))
		return
	}
	c.JSON(200, Result{
		Msg:  "ok",
		Data: id,
	})
}
