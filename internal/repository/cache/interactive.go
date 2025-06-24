package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/jym0818/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

var (
	//go:embed lua/interative_incr_cnt.lua
	luaIncrCnt string
)
var ErrKeyNotExist = redis.Nil

const (
	fieldReadCnt    = "read_cnt"
	fieldCollectCnt = "collect_cnt"
	fieldLikeCnt    = "like_cnt"
)

type InteractiveCache interface {

	// IncrReadCntIfPresent 如果在缓存中有对应的数据，就 +1
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error

	IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error
	// Get 查询缓存中数据
	// 事实上，这里 liked 和 collected 是不需要缓存的
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error
}

type interactiveCache struct {
	cmd redis.Cmdable
}

func NewinteractiveCache(cmd redis.Cmdable) InteractiveCache {
	return &interactiveCache{cmd: cmd}
}

func (cache *interactiveCache) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {

	// 拿到 key 对应的值里面的所有的 field
	data, err := cache.cmd.HGetAll(ctx, cache.key(biz, bizId)).Result()
	if err != nil {
		return domain.Interactive{}, err
	}

	if len(data) == 0 {
		// 缓存不存在，系统错误，比如说你的同事，手贱设置了缓存，但是忘记任何 fields
		return domain.Interactive{}, ErrKeyNotExist
	}

	// 理论上来说，这里不可能有 error
	collectCnt, _ := strconv.ParseInt(data[fieldCollectCnt], 10, 64)
	likeCnt, _ := strconv.ParseInt(data[fieldLikeCnt], 10, 64)
	readCnt, _ := strconv.ParseInt(data[fieldReadCnt], 10, 64)

	return domain.Interactive{
		CollectCnt: collectCnt,
		LikeCnt:    likeCnt,
		ReadCnt:    readCnt,
	}, err
}

func (cache *interactiveCache) Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error {
	key := cache.key(biz, bizId)
	err := cache.cmd.HMSet(ctx, key, fieldLikeCnt, intr.LikeCnt, fieldCollectCnt, intr.CollectCnt, fieldReadCnt, intr.ReadCnt).Err()
	if err != nil {
		return err
	}
	return cache.cmd.Expire(ctx, key, time.Minute*15).Err()
}

func (cache *interactiveCache) IncrCollectCntIfPresent(ctx context.Context,
	biz string, bizId int64) error {
	return cache.cmd.Eval(ctx, luaIncrCnt, []string{cache.key(biz, bizId), fieldCollectCnt}, 1).Err()
}

func (cache *interactiveCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return cache.cmd.Eval(ctx, luaIncrCnt, []string{cache.key(biz, bizId)}, fieldReadCnt, 1).Err()
}
func (cache *interactiveCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
func (cache *interactiveCache) IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return cache.cmd.Eval(ctx, luaIncrCnt, []string{cache.key(biz, bizId)}, fieldLikeCnt, 1).Err()
}

func (cache *interactiveCache) DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return cache.cmd.Eval(ctx, luaIncrCnt, []string{cache.key(biz, bizId)}, fieldLikeCnt, -1).Err()
}
