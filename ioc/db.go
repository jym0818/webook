package ioc

import (
	"github.com/jym0818/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(118.25.44.1:13316)/webook"))
	if err != nil {
		panic(err)
	}
	err = dao.InitDB(db)
	if err != nil {
		panic(err)
	}
	return db
}
