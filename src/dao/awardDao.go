package dao

import (
	"DrFinder/src/models"
	"github.com/jinzhu/gorm"
)

type AwardDao struct {
	engine *gorm.DB
}

func NewAwardDao(engine *gorm.DB) *AwardDao {
	return &AwardDao{engine:engine}
}

func (d *AwardDao) Add(award *models.Award) error {
	db := d.engine.Create(award)

	return db.Error
}

func (d *AwardDao) GetAwardByNpi(npi int) []*models.Award {
	var award []*models.Award

	db := d.engine.Where("npi = ?", npi).Find(&award)

	if db.Error != nil {
		return nil
	}

	return award
}
