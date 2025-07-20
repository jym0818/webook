package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

func main() {
	initViperV1()
	app := InitApp()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	err := app.server.Serve()
	log.Println(err)
	_ = app.server.Close()

}
func initViperV1() {
	cfgFile := pflag.String("config", "./config/config.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfgFile)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
