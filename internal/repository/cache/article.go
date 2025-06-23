package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jym0818/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error
	DelFirstPage(ctx context.Context, author int64) error
	Set(ctx context.Context, art domain.Article) error
	Get(ctx context.Context, id int64) (domain.Article, error)
	Del(ctx context.Context, id int64) error

	// SetPub 正常来说，创作者和读者的 Redis 集群要分开，因为读者是一个核心中的核心
	SetPub(ctx context.Context, article domain.Article) error
	GetPub(ctx context.Context, id int64) (domain.Article, error)
}

type articleCache struct {
	cmd redis.Cmdable
}

func (cache *articleCache) SetPub(ctx context.Context, art domain.Article) error {
	data, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return cache.cmd.Set(ctx, cache.readerArtKey(art.Id),
		data,
		// 设置长过期时间
		time.Minute*30).Err()
}

func (cache *articleCache) GetPub(ctx context.Context, id int64) (domain.Article, error) {
	// 可以直接使用 Bytes 方法来获得 []byte
	data, err := cache.cmd.Get(ctx, cache.readerArtKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var res domain.Article
	err = json.Unmarshal(data, &res)
	return res, err
}

func (cache *articleCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	data, err := cache.cmd.Get(ctx, cache.firstPageKey(uid)).Bytes()
	if err != nil {
		return nil, err
	}
	var arts []domain.Article
	err = json.Unmarshal(data, &arts)
	return arts, err
}

func (cache *articleCache) SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error {
	for i := range arts {
		// 只缓存摘要部分
		arts[i].Content = arts[i].Abstract()
	}
	data, err := json.Marshal(&arts)
	if err != nil {
		return err
	}
	return cache.cmd.Set(ctx, cache.firstPageKey(uid), data, time.Minute*10).Err()
}
func (cache *articleCache) DelFirstPage(ctx context.Context, author int64) error {
	return cache.cmd.Del(ctx, cache.firstPageKey(author)).Err()
}

func (cache *articleCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	// 可以直接使用 Bytes 方法来获得 []byte
	data, err := cache.cmd.Get(ctx, cache.authorArtKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var res domain.Article
	err = json.Unmarshal(data, &res)
	return res, err
}

func (cache *articleCache) Set(ctx context.Context, art domain.Article) error {
	data, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return cache.cmd.Set(ctx, cache.authorArtKey(art.Id), data, time.Minute).Err()
}

func (cache *articleCache) Del(ctx context.Context, id int64) error {
	return cache.cmd.Del(ctx, cache.authorArtKey(id)).Err()
}

func (cache *articleCache) authorArtKey(id int64) string {
	return fmt.Sprintf("article:author:%d", id)
}

func NewArticleCache(cmd redis.Cmdable) ArticleCache {
	return &articleCache{cmd: cmd}
}
func (cache *articleCache) firstPageKey(author int64) string {
	return fmt.Sprintf("article:first_page:%d", author)
}

// 读者端的缓存设置
func (cache *articleCache) readerArtKey(id int64) string {
	return fmt.Sprintf("article:reader:%d", id)
}
