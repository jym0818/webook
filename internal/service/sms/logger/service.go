package logger

import (
	"context"
	"github.com/jym/webook/internal/service/sms"
	"go.uber.org/zap"
)

type Service struct {
	svc sms.Service
}

func NewService(svc sms.Service) sms.Service {
	return &Service{svc: svc}
}
func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	zap.L().Debug("发送短息", zap.String("tpl", tpl), zap.Any("args", args))
	err := s.svc.Send(ctx, tpl, args, numbers...)
	if err != nil {
		zap.L().Debug("发送短信异常", zap.Error(err))
	}

	return err
}
