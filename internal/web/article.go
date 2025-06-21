package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/service"
	"go.uber.org/zap"
	"net/http"
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
	g.POST("/publish", h.Publish)
}
func (h *ArticleHandler) Edit(c *gin.Context) {

	var req ArticleReq
	if err := c.Bind(&req); err != nil {
		return
	}

	claims := c.MustGet("claims").(*UserClaims)

	//调用下一层

	id, err := h.svc.Save(c.Request.Context(), req.toDomain(claims.Uid))
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

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	claims := ctx.MustGet("claims").(*UserClaims)
	id, err := h.svc.Publish(ctx, req.toDomain(claims.Uid))

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 打日志？
		zap.L().Error("发布失败", zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})

}
func (req ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
