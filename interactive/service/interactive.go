package service

import (
	"context"
	"github.com/jym0818/webook/interactive/domain"
	"github.com/jym0818/webook/interactive/repository"
	"golang.org/x/sync/errgroup"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	// Like 点赞
	Like(ctx context.Context, biz string, bizId int64, uid int64) error
	// CancelLike 取消点赞
	CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error

	Collect(ctx context.Context, biz string, bizId, cid, uid int64) error
	Get(ctx context.Context, biz string, bizId, uid int64) (domain.Interactive, error)

	GetByIds(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository
}

func (svc *interactiveService) GetByIds(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func (svc *interactiveService) Get(ctx context.Context, biz string, bizId, uid int64) (domain.Interactive, error) {
	// 按照 repository 的语义(完成 domain.Interactive 的完整构造)，你这里拿到的就应该是包含全部字段的
	var (
		eg        errgroup.Group
		intr      domain.Interactive
		liked     bool
		collected bool
	)
	eg.Go(func() error {
		var err error
		intr, err = svc.repo.Get(ctx, biz, bizId)
		return err
	})
	eg.Go(func() error {
		var err error
		liked, err = svc.repo.Liked(ctx, biz, bizId, uid)
		return err
	})
	eg.Go(func() error {
		var err error
		liked, err = svc.repo.Collected(ctx, biz, bizId, uid)
		return err
	})
	err := eg.Wait()
	if err != nil {
		return domain.Interactive{}, err
	}
	intr.Liked = liked
	intr.Collected = collected
	return intr, err
}

func (svc *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return svc.repo.IncrReadCnt(ctx, biz, bizId)
}
func (svc *interactiveService) Like(ctx context.Context, biz string, bizId int64, uid int64) error {
	// 点赞
	return svc.repo.IncrLike(ctx, biz, bizId, uid)
}

func (svc *interactiveService) CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.DecrLike(ctx, biz, bizId, uid)
}
func (svc *interactiveService) Collect(ctx context.Context, biz string, bizId, cid, uid int64) error {
	// service 还叫做收藏
	// repository
	return svc.repo.AddCollectionItem(ctx, biz, bizId, cid, uid)
}

func NewinteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{
		repo: repo,
	}
}
