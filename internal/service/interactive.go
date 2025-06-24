package service

import (
	"context"
	"github.com/jym0818/webook/internal/repository"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	// Like 点赞
	Like(ctx context.Context, biz string, bizId int64, uid int64) error
	// CancelLike 取消点赞
	CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error
}

type interactiveService struct {
	repo repository.InteractiveRepository
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

func NewinteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{
		repo: repo,
	}
}
