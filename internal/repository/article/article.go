package article

import (
	"context"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/repository/dao/article"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	//存储并同步数据
	Sync(ctx context.Context, art domain.Article) (int64, error)
}

type CachedArticleRepository struct {
	dao article.ArticleDAO
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
	})
}

func (c *CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}
