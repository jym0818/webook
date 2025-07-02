package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jym0818/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type UserCache interface {
	Set(ctx context.Context, user domain.User) error
	Get(ctx context.Context, uid int64) (domain.User, error)
}

type userCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (cache *userCache) Set(ctx context.Context, user domain.User) error {
	u, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := cache.key(user.Id)
	return cache.cmd.Set(ctx, key, u, cache.expiration).Err()
}

func (cache *userCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func (cache *userCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	ctx = context.WithValue(ctx, "biz", "user")
	user, err := cache.cmd.Get(ctx, cache.key(uid)).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(user, &u)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func NewuserCache(cmd redis.Cmdable) UserCache {
	return &userCache{cmd: cmd, expiration: time.Minute * 15}
}
