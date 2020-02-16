package dao

import "github.com/jinzhu/gorm"

type FavorDao struct {
	engine *gorm.DB
}

func NewFavorDao(engine *gorm.DB) *FavorDao {
	return &FavorDao{engine: engine}
}

func (d *FavorDao) addFavor()  {
	
}
