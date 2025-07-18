package repository

import (
	"context"
	"github.com/jym0818/webook/interactive/domain"
	"github.com/jym0818/webook/interactive/repository/cache"
	"github.com/jym0818/webook/interactive/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrLike(ctx context.Context, biz string, bizId, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error

	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
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
func (repo *interactiveRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := repo.dao.GetLikeInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		// 你要吞掉
		return false, nil
	default:
		return false, err
	}
}

func (repo *interactiveRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := repo.dao.GetCollectionInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		// 你要吞掉
		return false, nil
	default:
		return false, err
	}
}

func (repo *interactiveRepository) AddCollectionItem(ctx context.Context, biz string, bizId, cid, uid int64) error {
	// 这个地方，你要不要考虑缓存收藏夹？
	// 以及收藏夹里面的内容
	// 用户会频繁访问他的收藏夹，那么你就应该缓存，不然你就不需要
	// 一个东西要不要缓存，你就看用户会不会频繁访问（反复访问）
	err := repo.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		Cid:   cid,
		Biz:   biz,
		BizId: bizId,
		Uid:   uid,
	})
	if err != nil {
		return err
	}
	// 收藏个数（有多少个人收藏了这个 biz + bizId)
	return repo.cache.IncrCollectCntIfPresent(ctx, biz, bizId)
}

func (repo *interactiveRepository) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	// 要从缓存拿出来阅读数，点赞数和收藏数
	intr, err := repo.cache.Get(ctx, biz, bizId)
	if err == nil {
		return intr, nil
	}

	daoIntr, err := repo.dao.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}
	intr = repo.toDomain(daoIntr)
	go func() {
		er := repo.cache.Set(ctx, biz, bizId, intr)
		// 记录日志
		if er != nil {

		}
	}()
	return intr, nil
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
func (repo *interactiveRepository) toDomain(intr dao.Interactive) domain.Interactive {
	return domain.Interactive{
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		ReadCnt:    intr.ReadCnt,
	}
}
