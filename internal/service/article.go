package service

import (
	"context"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/repository/article"
	"github.com/jym/webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, art domain.Article) error
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
	List(c context.Context, uid int64, limit int, offset int) ([]domain.Article, error)
	GetById(c context.Context, id int64) (domain.Article, error)
	GetPublishedById(c context.Context, id int64) (domain.Article, error)
}
type articleService struct {
	repo article.ArticleRepository
	//v1
	author article.ArticleAuthorRepository
	reader article.ArticleReaderRepository

	l logger.LoggerV1
}

func (a *articleService) GetPublishedById(c context.Context, id int64) (domain.Article, error) {
	return a.repo.GetPublishedById(c, id)
}

func (a *articleService) GetById(c context.Context, id int64) (domain.Article, error) {
	return a.repo.GetById(c, id)
}

func (a *articleService) List(c context.Context, uid int64, limit int, offset int) ([]domain.Article, error) {
	return a.repo.List(c, uid, limit, offset)
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func NewArticleServiceV1(author article.ArticleAuthorRepository, reader article.ArticleReaderRepository, l logger.LoggerV1) ArticleService {
	return &articleService{
		author: author,
		reader: reader,
		l:      l,
	}
}

func (a *articleService) Withdraw(ctx context.Context, art domain.Article) error {
	return a.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnPublished
	//如果有id  说明是修改   创建和修改共用一个接口
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)

}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return a.repo.Sync(ctx, art)
}

// 新建并发表
func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {

	var (
		id  = art.Id
		err error
	)
	if art.Id > 0 {
		err = a.author.Update(ctx, art)
	} else {
		id, err = a.author.Create(ctx, art)

	}
	if err != nil {
		return 0, err
	}
	//确保制作库和线上库的Id相等
	art.Id = id
	for i := 0; i < 3; i++ {
		id, err = a.reader.Save(ctx, art)
		if err == nil {
			break
		}
		a.l.Error("保存到线上库部分失败", logger.Int64("art_id", art.Id), logger.Error(err))

	}
	if err != nil {
		a.l.Error("重试失败", logger.Int64("art_id", art.Id), logger.Error(err))
	}
	return id, err
}
