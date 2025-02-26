package ioc

import (
	"github.com/jym/webook/internal/service/ratelimit"
	"github.com/jym/webook/internal/service/sms"
	"github.com/jym/webook/internal/service/sms/memory"
	ratelimit2 "github.com/jym/webook/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"time"
)

// 可以更改实现  memory或者tencent或者aliyun
func InitSMSService(cmd redis.Cmdable) sms.Service {
	svc := memory.NewService()
	limiter := ratelimit2.NewRedisSlidingWindowLimiter(cmd, time.Minute, 1000)
	return ratelimit.NewRatelimitSMSService(svc, limiter)
}
