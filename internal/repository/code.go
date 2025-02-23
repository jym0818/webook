package repository

import (
	"context"
	"github.com/jym/webook/internal/repository/cache"
)

var ErrSetCodeTooMany = cache.ErrSetCodeTooMany
var ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(cache *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: cache,
	}
}

func (repo *CodeRepository) Store(ctx context.Context, biz string, code string, phone string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}
func (repo *CodeRepository) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, inputCode)
}
