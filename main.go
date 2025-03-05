package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	initViperRemote()
	s := InitWebServer()
	s.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})
	s.Run(":8080")
}

func initViper() {
	viper.SetDefault("db.mysql.dsn", "root:root@tcp(127.0.0.1:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local")
	//文件名是config  但不包括文件后缀名
	viper.SetConfigName("dev")
	//文件格式是yaml
	viper.SetConfigType("yaml")
	//当前工作目录下的config子目录
	viper.AddConfigPath("./config")
	//读取配置到viper里面，或者你可以理解为加载到内存中
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}

func initViperV1() {
	viper.SetConfigFile("config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
func initViperRemote() {
	//etcd3.0之后的版本为etcd3
	//127.0.0.1:12379可以在配置文件中加载  也就是二次加载
	//配置在/webook下
	viper.SetConfigType("yaml")
	err := viper.AddRemoteProvider("etcd3", "127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func initViperV2() {
	//第一个参数 是参数名 	第二个参数是默认路径	第三个参数是注释
	cfile := pflag.String("config", "config/config.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
