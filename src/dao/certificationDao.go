package dao

import (
	"DrFinder/src/models"
	"github.com/jinzhu/gorm"
)

type CertificationDao struct {
	engine *gorm.DB
}

func NewCertificationDao(engine *gorm.DB) *CertificationDao {
	return &CertificationDao{engine:engine}
}

func (d *CertificationDao) Add(cer *models.Certification) error {
	db := d.engine.Create(cer)

	return db.Error
}

func (d *CertificationDao) GetCertificationByNpi(npi int) []models.Certification {
	var cer []models.Certification

	db := d.engine.Where("npi = ?", npi).Find(&cer)

	if db.Error != nil {
		return nil
	}

	return cer
}
