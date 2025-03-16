package article

import (
	"context"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/repository/cache"
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
	GetById(ctx context.Context, id int64) (domain.Article, error)
}

type CachedArticleRepository struct {
	dao   article.ArticleDAO
	cache cache.CacheArticle
}

func (c *CachedArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	res, err := c.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return c.toDomain(res), nil
}

func (c *CachedArticleRepository) List(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error) {
	//先去缓存中查找
	if offset == 0 && limit <= 100 {
		data, err := c.cache.GetFirstPage(ctx, uid)
		if err == nil {
			return data[:limit], nil
		}
	}

	//再去数据库
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
	//回写缓存.也可以异步
	//使用set
	err = c.cache.SetFirstPage(ctx, uid, arts)
	if err != nil {
		//记录日志
		//可以接受缓存失败
	}
	return arts, nil
}

func (c *CachedArticleRepository) SyncStatus(ctx context.Context, id int64, uid int64, status domain.ArticleStatus) error {
	return c.dao.SyncStatus(ctx, id, uid, status.ToUint8())
}

func NewCachedArticleRepository(dao article.ArticleDAO, cache cache.CacheArticle) ArticleRepository {
	return &CachedArticleRepository{
		dao:   dao,
		cache: cache,
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
	defer func() {
		//清空缓存
		c.cache.DelFirstPage(ctx, art.Author.Id)
	}()
	return c.dao.Sync(ctx, c.toEntity(art))
}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	defer func() {
		//清空缓存
		c.cache.DelFirstPage(ctx, art.Author.Id)
	}()
	return c.dao.Insert(ctx, article.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	})
}

func (c *CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	defer func() {
		//清空缓存
		c.cache.DelFirstPage(ctx, art.Author.Id)
	}()
	return c.dao.UpdateById(ctx, article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	})
}

func (c *CachedArticleRepository) toDomain(art article.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
	}
}
