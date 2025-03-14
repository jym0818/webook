package domain

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
}

type Author struct {
	Id   int64
	Name string
}

// 通常我们使用自定义类型 ，因为我们可以给这个类型定义自定义方法，用于业务
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

// 定义各种相关的业务方法，例如判断是否发布
func (s ArticleStatus) NonPublish() bool {
	return s != ArticleStatusUnPublished
}

// 对于简单的可以这么写，新增状态必须每次修改，复杂的使用下面的结构体
func (s ArticleStatus) String() string {
	switch s {
	case ArticleStatusUnKnown:
		return "UnKnown"
	case ArticleStatusUnPublished:
		return "UnPublished"
	case ArticleStatusPublished:
		return "Published"
	case ArticleStatusPrivate:
		return "Private"
	default:
		return "UnKnown"
	}

}

// 是否合法
func (s ArticleStatus) Valid() bool {
	return s.ToUint8() > 0
}
func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}

// 如果我们的状态很复杂，有很多行为，或者状态里面需要额外字段，可以定义为结构体，还是优先使用上面的方式
type ArticleStatusV1 struct {
	Val  uint8
	Name string
}

var (
	ArticleStatusV1UnKnown = ArticleStatusV1{Val: 0, Name: "UnKnown"}
)
