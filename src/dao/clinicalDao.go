package dao

import (
	"DrFinder/src/models"
	"github.com/jinzhu/gorm"
)

type ClinicalDao struct {
	engine *gorm.DB
}

func NewClinicalDao(engine *gorm.DB) *ClinicalDao {
	return &ClinicalDao{engine:engine}
}

func (d *ClinicalDao) Add(clinic *models.Clinical) error {
	db := d.engine.Create(clinic)
	return db.Error
}

func (d *ClinicalDao) GetClinicalByNpi(npi int) *models.Clinical {
	var clinical models.Clinical

	db := d.engine.Where("npi = ? ", npi).Find(&clinical)

	if db.Error != nil {
		return nil
	}

	return &clinical
}
