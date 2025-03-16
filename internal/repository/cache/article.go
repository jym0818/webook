package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jym/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheArticle interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, res []domain.Article) error
	DelFirstPage(ctx context.Context, uid int64) error
	Set(ctx context.Context, article domain.Article) error
	SetPub(ctx context.Context, art domain.Article) error
}
type RedisCacheArticle struct {
	cmd redis.Cmdable
}

func (r *RedisCacheArticle) SetPub(ctx context.Context, art domain.Article) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisCacheArticle) Set(ctx context.Context, article domain.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	//过期时间要短，你的预测效果越不好，时间越要短
	return r.cmd.Set(ctx, r.firstKey(article.Id), data, time.Second*10).Err()
}

func (r *RedisCacheArticle) firstKey(id int64) string {
	return fmt.Sprintf("article:%d", id)
}

func NewRedisCacheArticle(cmd redis.Cmdable) CacheArticle {
	return &RedisCacheArticle{
		cmd: cmd,
	}
}

func (r *RedisCacheArticle) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisCacheArticle) SetFirstPage(ctx context.Context, uid int64, res []domain.Article) error {
	for i := 0; i < len(res); i++ {
		res[i].Content = res[i].Abstract()
	}
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	return r.cmd.Set(ctx, r.key(uid), data, time.Minute*10).Err()
}

func (r *RedisCacheArticle) key(uid int64) string {
	return fmt.Sprintf("article:first_page:%d", uid)
}

func (r *RedisCacheArticle) DelFirstPage(ctx context.Context, uid int64) error {
	//TODO implement me
	panic("implement me")
}
