package failover

import (
	"context"
	"github.com/jym/webook/internal/service/sms"
	"sync/atomic"
)

type TimeoutFailoverSMSService struct {
	svcs []sms.Service
	idx  int32
	cnt  int32 //连续超时的个数
	//阈值  连续错误超过这个数字就要切换
	threshold int32
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt > t.threshold {
		//这里要切换 新的下标往后移动一位
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			//意味着我成功往后移动一位  idx+1
			atomic.StoreInt32(&t.cnt, 0)
		}
		//else 就是出现并发了 idx的值不正确了
		idx = atomic.LoadInt32(&t.idx)
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
		return err
	case nil:
		//连续状态打断了
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	default:
		return err
	}
}

func NewTimeoutFailoverSMSService() sms.Service {
	return &TimeoutFailoverSMSService{}
}
