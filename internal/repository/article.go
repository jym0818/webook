package repository

import (
	"context"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/repository/cache"
	"github.com/jym0818/webook/internal/repository/dao"
	"go.uber.org/zap"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, author int64, status domain.ArticleStatus) error
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetByID(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id int64) (domain.Article, error)
}

type articleRepository struct {
	dao   dao.ArticleDAO
	cache cache.ArticleCache

	userRepo UserRepository
}

func (repo *articleRepository) GetPublishedById(ctx context.Context, id int64) (domain.Article, error) {
	// 读取线上库数据，如果你的 Content 被你放过去了 OSS 上，你就要让前端去读 Content 字段
	art, err := repo.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	// 你在这边要组装 user 了，适合单体应用
	usr, err := repo.userRepo.FindById(ctx, art.AuthorId)
	if err != nil {
		return domain.Article{}, err
	}
	res := domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id:   usr.Id,
			Name: usr.Nickname,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
	return res, nil
}

func (repo *articleRepository) GetByID(ctx context.Context, id int64) (domain.Article, error) {
	//先查询缓存
	data, err := repo.cache.Get(ctx, id)
	if err == nil {
		return data, nil
	}
	//再查询数据库
	art, err := repo.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return repo.toDomain(art), nil
}

func (repo *articleRepository) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	if offset == 0 && limit <= 100 {
		data, err := repo.cache.GetFirstPage(ctx, uid)
		if err == nil {
			go func() {
				repo.preCache(ctx, data)
			}()
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
		err1 := repo.cache.SetFirstPage(ctx, uid, arts)
		if err1 != nil {
			//记录日志
		}
		repo.preCache(ctx, arts)

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
		er := repo.cache.SetPub(ctx, art)
		if er != nil {
			// 不需要特别关心
			// 比如说输出 WARN 日志
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

func NewarticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache, userRepo UserRepository) ArticleRepository {
	return &articleRepository{dao: dao, cache: cache, userRepo: userRepo}
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

func (repo *articleRepository) preCache(ctx context.Context, data []domain.Article) {
	if len(data) > 0 && len(data[0].Content) < 1024*1024 {
		err := repo.cache.Set(ctx, data[0])
		if err != nil {
			zap.L().Error("提前预加载缓存失败", zap.Error(err))
		}
	}
}
