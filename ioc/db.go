package ioc

import (
	"github.com/jym0818/webook/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
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
	err = dao.InitDB(db)
	if err != nil {
		panic(err)
	}
	return db
}
