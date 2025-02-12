package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

// 命名为Create 还是 Insert  个人偏好insert 因为更加贴近Mysql的操作
func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}

// 直接对应数据库表结构 与domain中user不是对应关系，可能会不同
// domain中的User是领域对象,是DDD中的entity或者聚合根，或者叫做BO
// 有些人叫做PO、entity，model，都一样
type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 设置为唯一索引
	Email    string `gorm:"unique"`
	Password string
	
	//创建时间 毫秒数 不使用time.Time  个人习惯 更加方便时区转换
	Ctime int64
	//更新时间 毫秒数
	Utime int64
}
