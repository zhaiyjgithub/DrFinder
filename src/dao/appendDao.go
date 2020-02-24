package dao

import (
	"DrFinder/src/models"
	"github.com/jinzhu/gorm"
)

type AppendDao struct {
	engine *gorm.DB
}

func NewAppendDao(engine *gorm.DB) *AppendDao {
	return &AppendDao{engine:engine}
}

func (d *AppendDao) AddAppend(append models.Append) error {
	db := d.engine.Create(append)

	return db.Error
}

func (d *AppendDao) GetAppends(postId int) []models.Append {
	var appends []models.Append
	d.engine.Where("post_id = ?", postId).Find(&appends)

	return appends
}
