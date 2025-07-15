package cache

import (
	"context"
	"encoding/json"
	"github.com/jym0818/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, arts []domain.Article) error
}

type RankingRedisCache struct {
	client redis.Cmdable
	key    string
}

func NewRankingRedisCache(client redis.Cmdable) RankingCache {
	return &RankingRedisCache{
		client: client,
		key:    "ranking",
	}

}
func (r *RankingRedisCache) Set(ctx context.Context, arts []domain.Article) error {
	// 你可以趁机，把 article 写到缓存里面 id => article
	for i := 0; i < len(arts); i++ {
		arts[i].Content = ""
	}
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	// 这个过期时间要稍微长一点，最好是超过计算热榜的时间（包含重试在内的时间）
	// 你甚至可以直接永不过期（建议使用----考虑计算热榜崩溃的兜底）
	return r.client.Set(ctx, r.key, val, time.Minute*10).Err()
}
