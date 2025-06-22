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
	g.POST("/withdraw", h.Withdraw)
	g.POST("list", h.List)
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
	id, err := h.svc.Publish(ctx.Request.Context(), req.toDomain(claims.Uid))

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

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	claims := ctx.MustGet("claims").(*UserClaims)
	err := h.svc.Withdraw(ctx.Request.Context(), domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 打日志？
		zap.L().Error("保存帖子失败", zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})

}

func (h *ArticleHandler) List(ctx *gin.Context) {
	var req ListReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	claims := ctx.MustGet("claims").(*UserClaims)
	res, err := h.svc.List(ctx.Request.Context(), claims.Uid, req.Limit, req.Offset)

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	var arts []ArticleVO
	for _, item := range res {
		arts = append(arts, ArticleVO{
			Id:       item.Id,
			Title:    item.Title,
			Status:   item.Status.ToUint8(),
			Ctime:    item.Ctime.Format("2006-01-02 15:04:05"),
			Utime:    item.Utime.Format("2006-01-02 15:04:05"),
			Abstract: item.Abstract(),
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: arts,
	})
}
