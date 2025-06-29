package web

import "github.com/jym0818/webook/internal/domain"

type ArticleVO struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	// 摘要
	Abstract string `json:"abstract"`
	// 内容
	Content string `json:"content"`
	// 注意一点，状态这个东西，可以是前端来处理，也可以是后端处理
	// 0 -> unknown -> 未知状态
	// 1 -> 未发表，手机 APP 这种涉及到发版的问题，那么后端来处理
	// 涉及到国际化，也是后端来处理
	Status uint8  `json:"status"`
	Author string `json:"author"`
	Ctime  string `json:"ctime"`
	Utime  string `json:"utime"`
	// 计数
	ReadCnt    int64 `json:"read_cnt"`
	LikeCnt    int64 `json:"like_cnt"`
	CollectCnt int64 `json:"collect_cnt"`

	// 我个人有没有收藏，有没有点赞
	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
}

type LikeReq struct {
	Id int64 `json:"id"`
	// 点赞和取消点赞，我都准备复用这个
	Like bool `json:"like"`
}

type ListReq struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
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

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
