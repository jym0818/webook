package article

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

func (dao *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (dao *GORMArticleDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	//建议使用这种方法，而不是使用零值忽略
	res := dao.db.WithContext(ctx).Model(&art).
		Where("id = ? AND author_id = ?", art.Id, art.AuthorId).Updates(map[string]interface{}{
		"title":   art.Title,
		"content": art.Content,
		"utime":   art.Utime,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("更新失败，可能是创作者非法id = %d", art.Id)
	}
	return nil
}

type Article struct {
	//如何设计索引
	//这是制作库
	//主要看WHERE
	//在帖子这里 什么样的查询场景
	//对于创作者来说，查看草稿箱，看到所有自己的文章
	//SELECT * FROM articles WHERE author_id = 123
	//单独查看某一篇 也就是查询id  主键  不需要加索引
	//同时产品经理会告诉你，要按照创建时间倒叙
	//所以最好在author_id和ctime上建立联合索引
	//SELECT * FROM articles WHERE author_id = 123 ORDER BY ctime DESC;
	Id       int64  `gorm:"primary_key;auto_increment"`
	Title    string `gorm:"type=varchar(1024)"`
	Content  string `gorm:"type=BLOB"`
	AuthorId int64  `gorm:"index"`
	Ctime    int64
	Utime    int64
}
