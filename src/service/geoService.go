package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type GeoService interface {
	Add(geo *models.Geo) error
	GetGeoInfoByNpi(npi int) *models.Geo
}

type geoService struct {
	dao *dao.GeoDao
}

func NewGeoService() GeoService {
	return &geoService{dao:dao.NewGeoDao(dataSource.InstanceMaster())}
}

func (s *geoService) Add(geo *models.Geo) error {
	return s.dao.Add(geo)
}

func (s *geoService) GetGeoInfoByNpi(npi int) *models.Geo {
	return s.dao.GetGeoInfoByNpi(npi)
}
