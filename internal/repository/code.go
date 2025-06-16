package repository

import (
	"context"
	"github.com/jym0818/webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
	ErrUnknownForCode         = cache.ErrUnknownForCode
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz string, phone string, code string) (bool, error)
}

type codeRepository struct {
	cache cache.CodeCache
}

func (repo *codeRepository) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, code)
}

func (repo *codeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func newCodeRepository(cache cache.CodeCache) CodeRepository {
	return &codeRepository{cache: cache}
}
