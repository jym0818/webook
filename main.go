package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	InitViper()
	server := InitServer()
	server.Run(":8080")
}

func InitViper() {

	cfgFile := pflag.String("config", "/config/dev.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfgFile)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
