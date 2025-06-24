package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/interative_incr_cnt.lua
	luaIncrCnt string
)

const (
	fieldReadCnt    = "read_cnt"
	fieldCollectCnt = "collect_cnt"
	fieldLikeCnt    = "like_cnt"
)

type InteractiveCache interface {

	// IncrReadCntIfPresent 如果在缓存中有对应的数据，就 +1
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
}

type interactiveCache struct {
	cmd redis.Cmdable
}

func NewinteractiveCache(cmd redis.Cmdable) InteractiveCache {
	return &interactiveCache{cmd: cmd}
}

func (cache *interactiveCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return cache.cmd.Eval(ctx, luaIncrCnt, []string{cache.key(biz, bizId)}, fieldReadCnt, 1).Err()
}
func (cache *interactiveCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
