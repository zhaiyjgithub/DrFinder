package dao

import (
	"DrFinder/src/models"
	"github.com/jinzhu/gorm"
)

type EducationDao struct {
	engine *gorm.DB
}

func NewEducationDao(engine *gorm.DB) *EducationDao {
	return &EducationDao{engine:engine}
}

func (d *EducationDao) Add(edu *models.Education) error {
	db := d.engine.Create(edu)

	return db.Error
}

func (d *EducationDao) GetEducationByNpi(npi int) []*models.Education {
	var edu []*models.Education
	db := d.engine.Where("npi = ?", npi).Find(&edu)
	if db.Error != nil {
		return nil
	}

	return edu
}


