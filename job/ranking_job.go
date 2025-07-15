package job

import (
	"context"
	"github.com/jym0818/webook/internal/service"
	"time"
)

type RankingJob struct {
	svc     service.RankingService
	timeout time.Duration
}

func NewRankingJob(svc service.RankingService, timeout time.Duration) *RankingJob {
	// 根据你的数据量来，如果要是七天内的帖子数量很多，你就要设置长一点
	return &RankingJob{
		svc:     svc,
		timeout: timeout,
	}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

// 按时间调度的，三分钟一次
func (r *RankingJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.svc.TopN(ctx)
}
