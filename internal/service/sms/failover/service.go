package failover

import (
	"context"
	"errors"
	"github.com/jym0818/webook/internal/service/sms"
	"sync/atomic"
)

type FailoverService struct {
	svcs []sms.Service
	idx  *atomic.Int64
}

func (s *FailoverService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := s.idx.Add(1)
	length := int64(len(s.svcs))
	for i := idx; i < idx+length; i++ {
		err := s.svcs[i%length].Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return nil
		default:
			//记录日志和监控
		}
	}
	return errors.New("failover service ")
}

func NewFailoverService(svcs []sms.Service) sms.Service {
	return &FailoverService{
		svcs: svcs,
		idx:  &atomic.Int64{},
	}
}
