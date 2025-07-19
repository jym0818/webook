package web

import (
	"github.com/gin-gonic/gin"
	intrv1 "github.com/jym0818/webook/api/proto/gen/intr/v1"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/service"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
)

type ArticleHandler struct {
	svc     service.ArticleService
	intrSvc intrv1.InteractiveServiceClient
	biz     string
}

func NewArticleHandler(svc service.ArticleService, intrSvc intrv1.InteractiveServiceClient) *ArticleHandler {
	return &ArticleHandler{
		svc:     svc,
		intrSvc: intrSvc,
		biz:     "article",
	}
}

func (h *ArticleHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/article")
	g.POST("/edit", h.Edit)
	g.POST("/publish", h.Publish)
	g.POST("/withdraw", h.Withdraw)
	g.POST("list", h.List)
	g.GET("/detail/:id", h.Detail)

	pub := g.Group("/pub")
	//pub.GET("/pub", a.PubList)
	pub.GET("/:id", h.PubDetail)

	pub.POST("/like", h.Like)

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

func (h *ArticleHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
	}

	claims := ctx.MustGet("claims").(*UserClaims)
	art, err := h.svc.GetById(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})

		return
	}
	// 这是不借助数据库查询来判定的方法
	if art.Author.Id != claims.Uid {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			// 也不需要告诉前端究竟发生了什么
			Msg: "输入有误",
		})
		// 如果公司有风控系统，这个时候就要上报这种非法访问的用户了。

		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVO{
			Id:    art.Id,
			Title: art.Title,
			// 不需要这个摘要信息
			//Abstract: art.Abstract(),
			Status:  art.Status.ToUint8(),
			Content: art.Content,
			// 这个是创作者看自己的文章列表，也不需要这个字段
			//Author: art.Author
			Ctime: art.Ctime.Format(time.DateTime),
			Utime: art.Utime.Format(time.DateTime),
		},
	})

}

func (h *ArticleHandler) PubDetail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})

		return
	}
	claims := ctx.MustGet("claims").(*UserClaims)
	var eg errgroup.Group
	var art domain.Article
	eg.Go(func() error {
		art, err = h.svc.GetPublishedById(ctx.Request.Context(), id, claims.Uid)
		return err
	})
	var intr *intrv1.GetResponse
	eg.Go(func() error {
		uc := ctx.MustGet("claims").(*UserClaims)
		intr, err = h.intrSvc.Get(ctx, &intrv1.GetRequest{
			Uid:   uc.Uid,
			Biz:   h.biz,
			BizId: id,
		})
		return err
	})
	err = eg.Wait()
	if err != nil {
		// 代表查询出错了
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	//增加阅读计数
	//go func() {
	//	er := h.intrSvc.IncrReadCnt(ctx, h.biz, art.Id)
	//	if er != nil {
	//		//记录日志
	//		zap.L().Error("记录日志失败 ", zap.Int64("文章Id", art.Id), zap.Error(er))
	//	}
	//}()

	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVO{
			Id:      art.Id,
			Title:   art.Title,
			Status:  art.Status.ToUint8(),
			Content: art.Content,
			// 要把作者信息带出去
			Author:     art.Author.Name,
			Ctime:      art.Ctime.Format(time.DateTime),
			Utime:      art.Utime.Format(time.DateTime),
			Liked:      intr.GetIntr().GetLiked(),
			Collected:  intr.GetIntr().GetCollected(),
			LikeCnt:    intr.GetIntr().GetLikeCnt(),
			ReadCnt:    intr.GetIntr().GetReadCnt(),
			CollectCnt: intr.GetIntr().GetCollectCnt(),
		},
	})
}

func (h *ArticleHandler) Like(ctx *gin.Context) {
	var req LikeReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	claims := ctx.MustGet("claims").(*UserClaims)
	var err error
	if req.Like {
		//h.biz, req.Id, claims.Uid
		_, err = h.intrSvc.Like(ctx.Request.Context(), &intrv1.LikeRequest{BizId: req.Id, Biz: h.biz, Uid: claims.Uid})
	} else {
		_, err = h.intrSvc.CancelLike(ctx.Request.Context(), &intrv1.CancelLikeRequest{BizId: req.Id, Biz: h.biz, Uid: claims.Uid})
	}
	if err != nil {
		zap.L().Error("点赞错误", zap.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})

}
