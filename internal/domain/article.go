package domain

import "time"

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Ctime   time.Time
	Utime   time.Time
	Status  ArticleStatus
}

type Author struct {
	Id   int64
	Name string
}

type ArticleStatus uint8

const (
	//未知状态   对于0值通常建议设置为无意义的值
	ArticleStatusUnKnown ArticleStatus = iota
	//未发表
	ArticleStatusUnPublished
	//已发表
	ArticleStatusPublished
	//仅自己可见
	ArticleStatusPrivate
)

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}
