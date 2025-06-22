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

func (a Article) Abstract() string {
	// 摘要我们取前几句。
	// 要考虑一个中文问题
	cs := []rune(a.Content)
	if len(cs) < 100 {
		return a.Content
	}
	// 英文怎么截取一个完整的单词，我的看法是……不需要纠结，就截断拉到
	// 词组、介词，往后找标点符号
	return string(cs[:100])
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
