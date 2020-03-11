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

func (d *GeoDao) GetNearByDoctorGeoInfo(lat float64, lng float64, page int, pageSize int) []models.GeoDistance {
	var geos []models.GeoDistance

	db := d.engine.Limit(pageSize).Offset(pageSize*(page - 1)).Order("distance asc").Raw("select npi, lat, lng, ACOS(SIN((? * 3.1415) / 180 ) *SIN((lat * 3.1415) / 180 ) +COS((? * 3.1415) / 180 ) * COS((lat * 3.1415) / 180 ) *COS((? * 3.1415) / 180 - (lng * 3.1415) / 180 ) ) * 6380 as distance  from geos where lat > (? - 1) and lat < (? + 1) and lng > (? - 1) and lng < (? + 1) ", lat, lat, lng, lat, lat, lng, lng).Scan(&geos)
	if db.Error != nil {
		return nil
	}

	return geos
}

//

func (d *GeoDao) GetDoctorGeoInfoByNpiList(lat float64, lng float64, npiList []int) []models.GeoDistance {
	var geos []models.GeoDistance
	db := d.engine.Raw("select npi, lat, lng, ACOS(SIN((? * 3.1415) / 180 ) *SIN((lat * 3.1415) / 180 ) +COS((? * 3.1415) / 180 ) * COS((lat * 3.1415) / 180 ) *COS((?* 3.1415) / 180 - (lng * 3.1415) / 180 ) ) * 6380 as distance from geos  where npi in (?)", lat,lng, lng, npiList).Scan(&geos)

	if db.Error != nil {
		return nil
	}

	return geos
}
