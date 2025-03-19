package web

import "github.com/jym/webook/internal/domain"

//VO   直接与前端打交道

type ArticleVo struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Abstract   string `json:"abstract"`
	Content    string `json:"content"`
	AuthorId   int64  `json:"author_id"`
	AuthorName string `json:"author_name"`
	//注意一点  状态这个东西，可以前端处理也可以后端处理
	//0 ---->Unkown->未知状态   这个转换可以自己绝对谁转换
	//如果APP设计到发版的问题  由后端处理
	//设计到国际化   也由后端处理
	Status uint8  `json:"status"`
	Ctime  string `json:"ctime"`
	Utime  string `json:"utime"`

	// 点赞之类的信息
	LikeCnt    int64 `json:"likeCnt"`
	CollectCnt int64 `json:"collectCnt"`
	ReadCnt    int64 `json:"readCnt"`

	// 个人是否点赞的信息
	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
}

type ListReq struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (req ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}
