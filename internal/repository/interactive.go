package repository

import (
	"context"
	"github.com/jym0818/webook/internal/repository/cache"
	"github.com/jym0818/webook/internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrLike(ctx context.Context, biz string, bizId, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId, uid int64) error
}
type interactiveRepository struct {
	cache cache.InteractiveCache
	dao   dao.InteractiveDAO
}

func NewinteractiveRepository(cache cache.InteractiveCache, dao dao.InteractiveDAO) InteractiveRepository {
	return &interactiveRepository{
		cache: cache,
		dao:   dao,
	}
}

func (repo *interactiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	// 要考虑缓存方案了
	// 这两个操作能不能换顺序？ —— 不能
	err := repo.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	//go func() {
	//	c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
	//}()
	//return err

	return repo.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}
func (repo *interactiveRepository) IncrLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	// 先插入点赞，然后更新点赞计数，更新缓存
	err := repo.dao.InsertLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	// 这种做法，你需要在 repository 层面上维持住事务
	//c.dao.IncrLikeCnt()
	return repo.cache.IncrLikeCntIfPresent(ctx, biz, bizId)
}

func (repo *interactiveRepository) DecrLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	err := repo.dao.DeleteLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return repo.cache.DecrLikeCntIfPresent(ctx, biz, bizId)
}
