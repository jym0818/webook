package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, article Article) (int64, error)
	Upsert(ctx context.Context, art PublishedArticle) error
}

type articleDAO struct {
	db *gorm.DB
}

func (dao *articleDAO) Sync(ctx context.Context, article Article) (int64, error) {
	var (
		id = article.Id
	)
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		txDAO := NewarticleDAO(tx)
		if id > 0 {
			err = txDAO.UpdateById(ctx, article)
		} else {
			id, err = txDAO.Insert(ctx, article)
		}
		if err != nil {
			return err
		}
		article.Id = id

		return txDAO.Upsert(ctx, PublishedArticle{article})

	})
	return id, err
}

func (dao *articleDAO) Upsert(ctx context.Context, art PublishedArticle) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	art.Ctime = now
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"utime":   art.Utime,
			"content": art.Content,
			"title":   art.Title,
		}),
	}).Create(&art).Error
}

func (dao *articleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err

}

func (dao *articleDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	res := dao.db.WithContext(ctx).Model(&art).Where("id = ? AND author_id = ?", art.Id, art.AuthorId).Updates(map[string]any{
		"utime":   now,
		"content": art.Content,
		"title":   art.Title,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("更新失败，可能是创作者非法 id %d, author_id %d",
			art.Id, art.AuthorId)
	}
	return res.Error
}

func NewarticleDAO(db *gorm.DB) ArticleDAO {
	return &articleDAO{db: db}
}

type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 长度 1024
	Title   string `gorm:"type=varchar(1024)"`
	Content string `gorm:"type=BLOB"`

	AuthorId int64 `gorm:"index"`

	Ctime int64
	Utime int64
}
type PublishedArticle struct {
	Article
}
