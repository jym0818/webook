package dao

import (
	"gorm.io/gorm"
)

func InitDB(db *gorm.DB) error {
	err := db.AutoMigrate(&Interactive{}, &UserLikeBiz{}, &Collection{}, &UserCollectionBiz{})
	return err
}
