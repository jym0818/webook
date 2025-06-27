package dao

import (
	"gorm.io/gorm"
)

func InitDB(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &Article{}, &PublishedArticle{}, &Interactive{}, &UserLikeBiz{}, &Collection{}, &UserCollectionBiz{})
	return err
}
