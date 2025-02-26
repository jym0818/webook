package ratelimit

import (
	"context"
	"fmt"
	"github.com/jym/webook/internal/service/sms"
	"github.com/jym/webook/pkg/ratelimit"
)

var errLimited = fmt.Errorf("出发了限流")

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

func (s RatelimitSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	//你在这里加上代码  新特性
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		return fmt.Errorf("短信服务是否出现了限流问题:%w", err)
	}
	if limited {
		return errLimited
	}

	err = s.svc.Send(ctx, tpl, args, numbers...)
	//你也可以在这里代码  新特性
	return err
}
