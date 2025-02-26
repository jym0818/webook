package ratelimit

import (
	"context"
	"fmt"
	"github.com/jym/webook/internal/service/sms"
	"github.com/jym/webook/pkg/ratelimit"
)

// 通过组合 默认实现了sms.Service接口
type RatelimitSMSServiceV1 struct {
	sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSServiceV1(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSServiceV1{
		Service: svc,
		limiter: limiter,
	}
}

func (s RatelimitSMSServiceV1) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	//你在这里加上代码  新特性
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		return fmt.Errorf("短信服务是否出现了限流问题:%w", err)
	}
	if limited {
		return errLimited
	}

	err = s.Send(ctx, tpl, args, numbers...)
	//你也可以在这里代码  新特性
	return err
}
