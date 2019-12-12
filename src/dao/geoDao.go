package dao

import (
	"DrFinder/src/models"
	"github.com/jinzhu/gorm"
)

type GeoDao struct {
	engine *gorm.DB
}

func NewGeoDao(engine *gorm.DB) *GeoDao {
	return &GeoDao{engine: engine}
}

func (d *GeoDao) Add(geo *models.Geo) error {
	db := d.engine.Create(geo)

	return db.Error
}

func (d *GeoDao) GetGeoInfoByNpi(npi int) *models.Geo {
	var geo = &models.Geo{}
	db := d.engine.Where("npi = ?", npi).Find(geo)

	if db.Error != nil {
		return nil
	}

	return geo
}


