package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jym/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExists = redis.Nil

type UserCache struct {
	//面向接口编程
	//我们就可以传入单机redis
	//也可以传入cluster的redis
	//也可以传入集群redis
	//不要将字段设置为 cmd *redis.Client 限制了扩展
	client redis.Cmdable
	//过期时间
	expiration time.Duration
}

// 依赖注入
// A用到了B  B一定是接口
// A用到了B  B一定是A的字段
// A用到了B A绝对不初始化B，而是外面注入
func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{client: client, expiration: time.Minute * 15}
}

// 即使没有数据  err也不能为nil 返回一个特定的err
// 只要err为nil  必然有数据
func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	//数据不存在返回redis.Nil
	str, err := cache.client.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User{}
	err = json.Unmarshal([]byte(str), &user)
	return user, err

}
func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(&u)
	if err != nil {
		return err
	}
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()

}
func (cache *UserCache) key(id int64) string {
	//user:info:123
	return fmt.Sprintf("user:info:%d", id)
}
