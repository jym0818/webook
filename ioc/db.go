package ioc

import (
	"github.com/jym/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		//我只会在初始化过程中panic
		//一旦panic goroutine就会结束
		//一旦初始化过程出错，应用就不要启动了，所以panic
		panic(err)
	}
	//初始化建表---------实际工作中不会使用这种方法
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
