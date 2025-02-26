package failover

import (
	"context"
	"errors"
	"github.com/jym/webook/internal/service/sms"
	"log"
)

type FailoverSMSService struct {
	svcs []sms.Service
}

func (f FailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		//发送成功
		if err == nil {
			return nil
		}
		//输出日志   要做好监控
		log.Println("failover service send error:", err)
	}
	return errors.New("全部服务都失败了")
}

func NewFailoverSMSService(svcs []sms.Service) sms.Service {
	return &FailoverSMSService{svcs: svcs}
}
