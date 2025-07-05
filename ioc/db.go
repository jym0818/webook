package ioc

import (
	"github.com/jym0818/webook/internal/repository/dao"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
	"time"
)

func InitDB() *gorm.DB {

	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config
	err := viper.UnmarshalKey("db.mysql", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "webook",
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{&prometheus.MySQL{
			VariableNames: []string{"Threads_running"},
		}},
	}))
	if err != nil {
		panic(err)
	}
	err = db.Use(tracing.NewPlugin(tracing.WithDBName("webook"),
		//不要记录metric  我们使用了prometheus
		tracing.WithoutMetrics(),
		//不要记录查询参数，安全需求线上不要记录
		tracing.WithoutQueryVariables(),
	))
	//统计查询时间
	summary := prometheus2.NewSummaryVec(prometheus2.SummaryOpts{
		Namespace:  "jym",
		Subsystem:  "webook",
		Name:       "gorm_query_time",
		Help:       "统计sql查询时间",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001, 0.999: 0.0001},
	}, []string{"type", "table"})
	prometheus2.MustRegister(summary)
	err = db.Callback().Create().Before("*").Register("prometheus_create_before", func(db *gorm.DB) {
		start := time.Now()
		db.Set("start_time", start)
	})
	if err != nil {
		panic(err)
	}
	err = db.Callback().Create().After("*").Register("prometheus_create_after", func(db *gorm.DB) {
		val, _ := db.Get("start_time")

		start, ok := val.(time.Time)
		if !ok {
			return
		}
		table := db.Statement.Table
		if table == "" {
			table = "unknown"
		}
		summary.WithLabelValues("create", table).Observe(float64(time.Since(start).Milliseconds()))
	})
	if err != nil {
		panic(err)
	}
	//查询
	err = db.Callback().Query().Before("*").Register("prometheus_query_before", func(db *gorm.DB) {
		start := time.Now()
		db.Set("start_time", start)
	})
	if err != nil {
		panic(err)
	}
	err = db.Callback().Query().After("*").Register("prometheus_query_after", func(db *gorm.DB) {
		val, _ := db.Get("start_time")

		start, ok := val.(time.Time)
		if !ok {
			return
		}
		table := db.Statement.Table
		if table == "" {
			table = "unknown"
		}
		summary.WithLabelValues("query", table).Observe(float64(time.Since(start).Milliseconds()))
	})
	if err != nil {
		panic(err)
	}

	err = dao.InitDB(db)
	if err != nil {
		panic(err)
	}
	return db
}
