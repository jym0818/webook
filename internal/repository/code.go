package repository

import (
	"context"
	"github.com/jym/webook/internal/repository/cache"
)

var ErrSetCodeTooMany = cache.ErrSetCodeTooMany
var ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes

type CodeRepository interface {
	Store(ctx context.Context, biz string, code string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type CacheCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(cache cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		cache: cache,
	}
}

func (repo *CacheCodeRepository) Store(ctx context.Context, biz string, code string, phone string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}
func (repo *CacheCodeRepository) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, inputCode)
}
