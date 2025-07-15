package ioc

import (
	"github.com/jym0818/webook/internal/service"
	"github.com/jym0818/webook/job"
	"github.com/robfig/cron/v3"
	"time"
)

func InitRankingJob(svc service.RankingService) job.Job {
	return job.NewRankingJob(svc, time.Second*30)
}

func InitCronJob(ranking job.Job) *cron.Cron {
	res := cron.New(cron.WithSeconds())
	cbd := job.NewCronJobBuilder()
	// 这里每三分钟一次
	_, err := res.AddJob("0 */3 * * * ?", cbd.Build(ranking))
	if err != nil {
		panic(err)
	}
	return res
}
