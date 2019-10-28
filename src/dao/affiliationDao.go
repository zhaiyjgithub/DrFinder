package dao

import (
	"DrFinder/src/models"
	"github.com/jinzhu/gorm"
)

type AffiliationDao struct {
	engine *gorm.DB
}

func NewAffiliationDao(engine *gorm.DB) *AffiliationDao {
	return &AffiliationDao{
		engine: engine,
	}
}

func (d *AffiliationDao) Add(affiliation *models.Affiliation) error {
	db := d.engine.Create(affiliation)

	return db.Error
}

func (d *AffiliationDao) GetAffiliationByNpi(npi int) *models.Affiliation {
	var affiliation models.Affiliation
	db := d.engine.Where("npi = ?", npi).Find(&affiliation)
	
	if db.Error != nil {
		return nil
	}
	
	return &affiliation
}

func (d *AffiliationDao) GetAll(page int, pageSize int) error {
	var affiliations []models.Affiliation

	db := d.engine.Limit(pageSize).Offset((page - 1)* pageSize).Find(&affiliations)

	return db.Error
}