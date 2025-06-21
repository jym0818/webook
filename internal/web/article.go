package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/service"
	"go.uber.org/zap"
)

type ArticleHandler struct {
	svc service.ArticleService
}

func NewArticleHandler(svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}

func (h *ArticleHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/article")
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

	claims := c.MustGet("claims").(*UserClaims)

	//调用下一层

	id, err := h.svc.Save(c.Request.Context(), domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author:  domain.Author{Id: claims.Uid},
	})
	if err != nil {
		c.JSON(500, Result{Msg: "保存帖子失败"})
		zap.L().Error("保存帖子失败", zap.Error(err))
		return
	}
	c.JSON(200, Result{
		Msg:  "ok",
		Data: id,
	})
}
