package repository

import (
	"context"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/repository/cache"
	"github.com/jym0818/webook/internal/repository/dao"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, author int64, status domain.ArticleStatus) error
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
}

type articleRepository struct {
	dao   dao.ArticleDAO
	cache cache.ArticleCache
}

func (repo *articleRepository) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	if offset == 0 && limit <= 100 {
		data, err := repo.cache.GetFirstPage(ctx, uid)
		if err == nil {

			return data[:limit], nil
		}
	}
	res, err := repo.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	arts := []domain.Article{}
	for _, v := range res {
		arts = append(arts, repo.toDomain(v))
	}
	//回写缓存
	go func() {
		err = repo.cache.SetFirstPage(ctx, uid, arts)
		if err != nil {
			//记录日志
		}
	}()
	return arts, nil

}

func (repo *articleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := repo.dao.Sync(ctx, repo.toEntity(art))
	if err == nil {
		err1 := repo.cache.DelFirstPage(ctx, art.Author.Id)
		if err1 != nil {
			//不太关心，记录日志
		}
	}
	return id, err
}
func (repo *articleRepository) SyncStatus(ctx context.Context, id int64, author int64, status domain.ArticleStatus) error {
	return repo.dao.SyncStatus(ctx, id, author, uint8(status))
}
func (repo *articleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	defer func() {
		// 清空缓存
		repo.cache.DelFirstPage(ctx, art.Author.Id)
	}()
	return repo.dao.Insert(ctx, repo.toEntity(art))
}

func (repo *articleRepository) Update(ctx context.Context, art domain.Article) error {
	defer func() {
		// 清空缓存
		repo.cache.DelFirstPage(ctx, art.Author.Id)
	}()
	return repo.dao.UpdateById(ctx, repo.toEntity(art))
}

func NewarticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache) ArticleRepository {
	return &articleRepository{dao: dao, cache: cache}
}

func (repo *articleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Ctime:    art.Ctime.UnixMilli(),
		Utime:    art.Utime.UnixMilli(),
		Status:   art.Status.ToUint8(),
	}
}

func (repo *articleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author:  domain.Author{Id: art.AuthorId},
		Ctime:   time.UnixMilli(art.Ctime),
		Utime:   time.UnixMilli(art.Utime),
		Status:  domain.ArticleStatus(art.Status),
	}
}
