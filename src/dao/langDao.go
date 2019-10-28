package dao

import (
	"DrFinder/src/models"
	"github.com/jinzhu/gorm"
)

type LangDao struct {
	engine *gorm.DB
}

func NewLangDao(engine *gorm.DB) *LangDao {
	return &LangDao{engine:engine}
}

func (d *LangDao) Add(clinic *models.Lang) error {
	db := d.engine.Create(clinic)
	return db.Error
}

func (d *LangDao) GetLangByNpi(npi int) *models.Lang {
	var Lang models.Lang

	db := d.engine.Where("npi = ? ", npi).Find(&Lang)

	if db.Error != nil {
		return nil
	}

	return &Lang
}

