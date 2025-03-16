package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/service"
	ijwt "github.com/jym/webook/internal/web/jwt"
	"github.com/jym/webook/pkg/logger"
	"net/http"
	"strconv"
	"time"
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
	g.POST("/publish", h.Publish)
	g.POST("/withdraw", h.Withdraw)
	g.POST("/list", h.List)
	g.GET("/detail/:id", h.Detail)
	g.GET("/pub/:id", h.PubDetail)
}

func (h *ArticleHandler) Withdraw(c *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		c.JSON(200, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	claims := c.MustGet("claims").(ijwt.UserClaims)
	err := h.svc.Withdraw(c, domain.Article{Id: req.Id, Author: domain.Author{Id: claims.Uid}})
	if err != nil {
		c.JSON(200, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("修改状态失败", logger.Error(err))
		return
	}
	c.JSON(200, Result{
		Msg: "ok",
	})
}

func (h *ArticleHandler) Publish(c *gin.Context) {
	var req ArticleReq

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

	id, err := h.svc.Publish(c, req.toDomain(claims.Uid))
	if err != nil {
		c.JSON(200, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("发表帖子失败", logger.Error(err))
		return
	}
	c.JSON(200, Result{
		Msg:  "ok",
		Data: id,
	})

}

func (h *ArticleHandler) Edit(c *gin.Context) {

	var req ArticleReq
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
	id, err := h.svc.Save(c, req.toDomain(claims.Uid))
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

func (h *ArticleHandler) List(c *gin.Context) {
	var req ListReq
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
	res, err := h.svc.List(c, claims.Uid, req.Limit, req.Offset)
	if err != nil {
		c.JSON(200, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	//返回的数据  列表页不显示全文，而是一段摘要
	//简单的摘要就是前几句话
	//强大的摘要就是AI生成的
	var resp []ArticleVo
	for _, v := range res {
		a := ArticleVo{
			Id:       v.Id,
			Title:    v.Title,
			Abstract: v.Abstract(),
			Status:   v.Status.ToUint8(),
			Ctime:    v.Ctime.Format(time.DateTime),
			Utime:    v.Utime.Format(time.DateTime),
		}
		resp = append(resp, a)
	}

	c.JSON(200, Result{
		Msg:  "ok",
		Data: resp,
	})
}

func (h *ArticleHandler) Detail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})
		return
	}
	usr, _ := c.MustGet("user").(*ijwt.UserClaims)
	art, err := h.svc.GetById(c, id)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})

		return
	}
	// 这是不借助数据库查询来判定的方法
	if art.Author.Id != usr.Uid {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			// 也不需要告诉前端究竟发生了什么
			Msg: "输入有误",
		})
		// 如果公司有风控系统，这个时候就要上报这种非法访问的用户了。
		h.l.Error("非法访问文章，创作者 ID 不匹配",
			logger.Int64("uid", usr.Uid))
		return
	}
	c.JSON(http.StatusOK, Result{
		Data: ArticleVo{
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

func (h *ArticleHandler) PubDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})
		h.l.Error("前端输入的 ID 不对", logger.Error(err))
		return
	}

	art, err := h.svc.GetPublishedById(c, id)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("获得文章信息失败", logger.Error(err))
		return
	}
	c.JSON(http.StatusOK, Result{
		Data: ArticleVo{
			Id:      art.Id,
			Title:   art.Title,
			Status:  art.Status.ToUint8(),
			Content: art.Content,
			// 要把作者信息带出去
			AuthorName: art.Author.Name,
			Ctime:      art.Ctime.Format(time.DateTime),
			Utime:      art.Utime.Format(time.DateTime),
		},
	})
}
