package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

// ErrUserDuplicateEmail 这个算是 user 专属的
var ErrUserDuplicateEmail = errors.New("邮件冲突")

// 继续一层一层暴露出去知道repository层，要在service层判断返回的err是不是这个错误
var ErrUserNotFound = gorm.ErrRecordNotFound

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
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		//判断是否冲突
		const uniqueIndexErrNo uint16 = 1062
		if me.Number == uniqueIndexErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	//无需检查错误，找不到会返回ErrRecordNotFound和空结构体
	return u, err
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
