package repository

import (
	"context"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}
type CachedRankingRepository struct {
	redis cache.RankingCache
	local *cache.RankingLocalCache
}

func NewCachedRankingRepository(redis cache.RankingCache, local *cache.RankingLocalCache) RankingRepository {
	return &CachedRankingRepository{redis: redis, local: local}
}
func (c *CachedRankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	_ = c.local.Set(ctx, arts)
	return c.redis.Set(ctx, arts)
}
func (c *CachedRankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	data, err := c.local.Get(ctx)
	if err == nil {
		return data, nil
	}
	data, err = c.redis.Get(ctx)
	if err == nil {
		c.local.Set(ctx, data)
	} else {
		return c.local.ForceGet(ctx)
	}
	return data, err
}
