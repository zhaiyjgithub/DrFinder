package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type GeoService interface {
	Add(geo *models.Geo) error
	GetGeoInfoByNpi(npi int) *models.Geo
	GetNearByDoctorGeoInfo(lat float64, lng float64, page int, pageSize int) []models.GeoDistance
	GetDoctorGeoInfoByNpiList(lat float64, lng float64, npiList []int) []models.GeoDistance
	GetUnInitList(page int, pageSize int) []models.Geo
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

func (s *geoService) GetNearByDoctorGeoInfo(lat float64, lng float64, page int, pageSize int) []models.GeoDistance  {
	return s.dao.GetNearByDoctorGeoInfo(lat, lng, page, pageSize)
}

func (s *geoService) GetDoctorGeoInfoByNpiList(lat float64, lng float64, npiList []int) []models.GeoDistance  {
	return s.dao.GetDoctorGeoInfoByNpiList(lat, lng, npiList)
}

func (s *geoService) GetUnInitList(page int, pageSize int) []models.Geo  {
	return s.dao.GetUnInitList(page, pageSize)
}