package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type DoctorElasticService interface {
	AddDoctor(doctor *models.Doctor, lat float64, lon float64) error
	QueryDoctor(
		doctorName string,
		specialty string,
		gender int,
		state string,
		city string,
		address string,
		zipCode int,
		page int,
		pageSize int,
	) []int
	QueryNewByDoctor(
		lat float64,
		lon float64,
		distance string,
		page int,
		pageSize int,
	) []int
}

func NewDoctorElasticService() DoctorElasticService {
	return &doctorElasticService{dao: dao.NewDoctorElasticDao(dataSource.InstanceElasticSearchClient())}
}

type doctorElasticService struct {
	dao *dao.DoctorElasticDao
}

func (s *doctorElasticService) AddDoctor(doctor *models.Doctor, lat float64, lon float64) error {
	return s.dao.AddDoctor(doctor, lat, lon)
}

func (s *doctorElasticService) QueryDoctor(
	doctorName string,
	specialty string,
	gender int,
	state string,
	city string,
	address string,
	zipCode int,
	page int,
	pageSize int,
) []int  {
	return s.dao.QueryDoctor(doctorName, specialty, gender, state, city, address, zipCode, page, pageSize)
}

func (s *doctorElasticService) QueryNewByDoctor(lat float64,
	lon float64,
	distance string,
	page int,
	pageSize int,
) []int {
	return s.dao.QueryNewByDoctor(lat, lon, distance, page, pageSize)
}


