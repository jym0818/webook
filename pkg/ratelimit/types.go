package ratelimit

import "context"

type Limiter interface {
	//Limit限流对象  key是限流对象
	//true代表限流  error返回错误
	Limit(ctx context.Context, key string) (bool, error)
}
