package job

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type CronJobBuilder struct {
	p      *prometheus.SummaryVec
	tracer trace.Tracer
}

func NewCronJobBuilder() *CronJobBuilder {
	p := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "jym",
		Subsystem: "webook",
		Help:      "统计定时任务的执行情况",
		Name:      "cron_job",
	}, []string{"name", "success"})
	prometheus.MustRegister(p)
	return &CronJobBuilder{
		p:      p,
		tracer: otel.GetTracerProvider().Tracer("webook/internal/job"),
	}
}
func (b *CronJobBuilder) Build(job Job) cron.Job {
	name := job.Name()
	return cronJobFuncAdapter(func() error {
		_, span := b.tracer.Start(context.Background(), name)
		defer span.End()
		start := time.Now()
		zap.L().Info("任务开始", zap.String("job", name))
		var success bool
		defer func() {
			zap.L().Info("任务结束", zap.String("job", name))
			duration := time.Since(start).Milliseconds()
			b.p.WithLabelValues(name, strconv.FormatBool(success)).Observe(float64(duration))
		}()
		err := job.Run()
		success = err == nil
		if err != nil {
			span.RecordError(err)
			zap.L().Error("运行任务失败", zap.Error(err), zap.String("job", name))
		}
		return nil
	})
}

type cronJobFuncAdapter func() error

func (c cronJobFuncAdapter) Run() {
	_ = c()
}
