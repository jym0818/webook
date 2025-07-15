package repository

import (
	"context"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
}
type CachedRankingRepository struct {
	redis cache.RankingCache
}

func NewCachedRankingRepository(redis cache.RankingCache) RankingRepository {
	return &CachedRankingRepository{redis: redis}
}
func (c *CachedRankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	return c.redis.Set(ctx, arts)
}
