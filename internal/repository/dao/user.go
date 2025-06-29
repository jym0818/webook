package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var ErrUserDuplicateEmail = errors.New("邮件冲突")
var ErrUserNotFound = gorm.ErrRecordNotFound

type UserDAO interface {
	Insert(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, uid int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByWechat(ctx context.Context, openId string) (User, error)
}

type userDAO struct {
	db *gorm.DB
}

func (dao *userDAO) FindByWechat(ctx context.Context, openID string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("openID = ?", openID).First(&u).Error
	//无需检查错误，找不到会返回ErrRecordNotFound和空结构体
	return u, err
}

func (dao *userDAO) FindById(ctx context.Context, uid int64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("id= ?", uid).First(&user).Error
	return user, err
}

func (dao *userDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email= ?", email).First(&user).Error
	return user, err
}

func (dao *userDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("phone= ?", phone).First(&user).Error
	return user, err
}
func NewuserDAO(db *gorm.DB) UserDAO {
	return &userDAO{db: db}
}

func (dao *userDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().Unix()
	user.Ctime = now
	user.Utime = now
	err := dao.db.WithContext(ctx).Create(&user).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const uniqueIndexErrNo uint16 = 1062
		if me.Number == uniqueIndexErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

type User struct {
	Id            int64          `gorm:"primaryKey,autoIncrement"`
	Email         sql.NullString `gorm:"unique"`
	Phone         sql.NullString `gorm:"unique"`
	Password      string
	Ctime         int64
	Utime         int64
	WechatUnionID sql.NullString
	WechatOpenID  sql.NullString `gorm:"unique"`
	//openID 一定是唯一的
	Nickname string
}
