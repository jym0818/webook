package domain

import "time"

type User struct {
	Id         int64
	Email      string
	Password   string
	Phone      string
	Ctime      time.Time
	Utime      time.Time
	WechatInfo WechatInfo
	Nickname   string
}
