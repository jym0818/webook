package repository

import (
	"context"
	"github.com/jym0818/webook/internal/repository/cache"
	"github.com/jym0818/webook/internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
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
