package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	InitViper()
	InitLogger()
	app := InitServer()
	for _, consumer := range app.consumers {
		err := consumer.Start()
		if err != nil {
			panic(err)
		}
	}
	app.web.Run(":8080")
}

func InitViper() {

	cfgFile := pflag.String("config", "./config/dev.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfgFile)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
