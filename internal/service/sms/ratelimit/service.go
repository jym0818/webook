package ratelimit

import (
	"context"
	"fmt"
	"github.com/jym0818/webook/internal/service/sms"
	"github.com/jym0818/webook/pkg/ratelimit"
)

var ErrLimited = fmt.Errorf("触发了限流")

type Service struct {
	sms     sms.Service
	limiter ratelimit.Limiter
}

func NewService(sms sms.Service, l ratelimit.Limiter) sms.Service {
	return &Service{
		sms:     sms,
		limiter: l,
	}
}
func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		return ErrLimited
	}
	if limited {
		return ErrLimited
	}
	return s.sms.Send(ctx, tpl, args, numbers...)
}
