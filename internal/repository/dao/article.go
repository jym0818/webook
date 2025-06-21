package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
}

type articleDAO struct {
	db *gorm.DB
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
