package article

import (
	"context"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/repository/dao/article"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	//存储并同步数据
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, uid int64, status domain.ArticleStatus) error
	List(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error)
}

type CachedArticleRepository struct {
	dao article.ArticleDAO
}

func (c *CachedArticleRepository) List(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error) {
	res, err := c.dao.GetByAuthor(ctx, uid, limit, offset)
	if err != nil {
		return nil, err
	}
	var arts []domain.Article
	for _, v := range res {
		arts = append(arts, domain.Article{
			Id:      v.Id,
			Title:   v.Title,
			Content: v.Content,
			Author: domain.Author{
				Id: v.AuthorId,
			},
			Ctime:  time.UnixMilli(v.Ctime),
			Utime:  time.UnixMilli(v.Utime),
			Status: domain.ArticleStatus(v.Status),
		})
	}
	return arts, nil
}

func (c *CachedArticleRepository) SyncStatus(ctx context.Context, id int64, uid int64, status domain.ArticleStatus) error {
	return c.dao.SyncStatus(ctx, id, uid, status.ToUint8())
}

func NewCachedArticleRepository(dao article.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

func (c *CachedArticleRepository) toEntity(article2 domain.Article) article.Article {
	return article.Article{
		Id:       article2.Id,
		Title:    article2.Title,
		Content:  article2.Content,
		AuthorId: article2.Author.Id,
		Status:   article2.Status.ToUint8(),
	}
}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Sync(ctx, c.toEntity(art))
}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, article.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	})
}

func (c *CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	})
}
