package article

import (
	"context"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/repository/dao/article"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	//存储并同步数据
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncV2(ctx context.Context, art domain.Article) (int64, error)
}

type CachedArticleRepository struct {
	dao article.ArticleDAO

	//v1  操作两个dao
	readerDAO article.ReaderDAO
	authorDAO article.AuthorDAO

	//在repository实现事务  必须依赖db
	db *gorm.DB
}

func NewCachedArticleRepository(dao article.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

// syncV2尝试在repository层面实现事务
// 确保保存到制作库和线上库同时成功  或者同时失败
// 如果过程中发生panic，不会回滚也不会提交  事务就会一直挂在数据库上，事务的结束会看数据库的默认配置（较长时间）
// 所以使用defer
func (c *CachedArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	//开启事务
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer tx.Rollback()
	//利用tx来构建
	author := dao.NewGormArticleAuthorDao(tx)
	reader := dao.NewGormArticleReaderDAO(tx)

	var (
		id  = art.Id
		err error
	)
	artn := c.toEntity(art)
	//先保存到制作库
	if art.Id > 0 {
		err = author.UpdateById(ctx, artn)
	} else {
		id, err = author.Insert(ctx, artn)
	}
	if err != nil {

		return id, err
	}
	//再保存到线上库
	//考虑到线上库可能有 可能没有
	//需要insert或者update
	err = reader.Upsert(ctx, artn)
	tx.Commit()
	return id, err

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
	var (
		id  = art.Id
		err error
	)
	artn := c.toEntity(art)
	//先保存到制作库
	if art.Id > 0 {
		err = c.authorDAO.UpdateById(ctx, artn)
	} else {
		id, err = c.authorDAO.Insert(ctx, artn)
	}
	if err != nil {
		return id, err
	}
	//再保存到线上库
	//考虑到线上库可能有 可能没有
	//需要insert或者update
	err = c.readerDAO.Upsert(ctx, artn)
	return id, err
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
