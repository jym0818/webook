package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

// ErrUserDuplicateEmail 这个算是 user 专属的
var ErrUserDuplicate = errors.New("邮件或者手机号冲突")

// 继续一层一层暴露出去知道repository层，要在service层判断返回的err是不是这个错误
var ErrUserNotFound = gorm.ErrRecordNotFound

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByWechat(ctx context.Context, openID string) (User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

// 命名为Create 还是 Insert  个人偏好insert 因为更加贴近Mysql的操作
func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		//判断是否冲突
		const uniqueIndexErrNo uint16 = 1062
		if me.Number == uniqueIndexErrNo {
			return ErrUserDuplicate
		}
	}
	return err
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	//无需检查错误，找不到会返回ErrRecordNotFound和空结构体
	return u, err
}

func (dao *GORMUserDAO) FindByWechat(ctx context.Context, openID string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("openID = ?", openID).First(&u).Error
	//无需检查错误，找不到会返回ErrRecordNotFound和空结构体
	return u, err
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	//无需检查错误，找不到会返回ErrRecordNotFound和空结构体
	return u, err
}

func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}

// 直接对应数据库表结构 与domain中user不是对应关系，可能会不同
// domain中的User是领域对象,是DDD中的entity或者聚合根，或者叫做BO
// 有些人叫做PO、entity，model，都一样
type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 设置为唯一索引
	Email    sql.NullString `gorm:"unique"`
	Password string

	Phone sql.NullString `gorm:"unique"`
	//创建时间 毫秒数 不使用time.Time  个人习惯 更加方便时区转换
	Ctime int64
	//更新时间 毫秒数
	Utime         int64
	WechatUnionID sql.NullString
	WechatOpenID  sql.NullString `gorm:"unique"`
}
