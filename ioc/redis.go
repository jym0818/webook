package ioc

import (
	"github.com/jym0818/webook/pkg/redisx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	if err := viper.UnmarshalKey("redis", &cfg); err != nil {
		panic(err)
	}
	cmd := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})
	cmd.AddHook(redisx.NewPrometheusHook(prometheus.SummaryOpts{
		Namespace:  "jym",
		Subsystem:  "webook",
		Name:       "redis",
		Help:       "监控reids",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}))
	return cmd
}
