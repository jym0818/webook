package ioc

import (
	"github.com/jym0818/webook/internal/service/sms"
	"github.com/jym0818/webook/internal/service/sms/memory"
	"github.com/jym0818/webook/internal/service/sms/ratelimit"
	ratelimit2 "github.com/jym0818/webook/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitSMS(cmd redis.Cmdable) sms.Service {
	smsSvc := memory.NewService()
	l := ratelimit2.NewRedisSlideWindow(cmd, time.Second, 2000)
	ratelimitSvc := ratelimit.NewService(smsSvc, l)
	return ratelimitSvc
}
