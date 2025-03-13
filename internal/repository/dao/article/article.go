package article

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	Upsert(ctx context.Context, art PublishArticle) error
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

func (dao *GORMArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	//先操作制作表  再操作 线上表  同一个库
	//采用闭包的形式
	//gorm帮我们处理事务的生命周期 不需要操心commit或者rollback
	//注意一条sql预计不需要开启事务
	var (
		id  = art.Id
		err error
	)
	err = dao.db.Transaction(func(tx *gorm.DB) error {
		var err error
		txDAO := NewGORMArticleDAO(tx)
		if id > 0 {
			err = txDAO.UpdateById(ctx, art)
		} else {
			id, err = txDAO.Insert(ctx, art)
		}
		if err != nil {
			return err
		}
		//操作线上库
		return txDAO.Upsert(ctx, PublishArticle{Article: art})
	})
	return id, err
}

func (dao *GORMArticleDAO) Upsert(ctx context.Context, art PublishArticle) error {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	//OnConflict的意思是数据冲突了
	err := dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		//哪些列冲突
		//Columns: []clause.Column{{Name: "id"}},
		//意思是数据冲突了 什么也不干
		//DoNothing: true,
		//数据冲突了并且符合where条件就会触发
		//Where: clause.Where{}
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
			"utime":   art.Utime,
		}),
	}).Create(&art).Error
	//最终的语句是
	//INSERT  xxx ON DUPLICATE KEY UPDATE xxx
	return err
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
