package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type DoctorService interface {
	Add(doctor *models.Doctor) bool
	GetDoctorById(id int) *models.Doctor
	GetDoctorBySpecialty(specialty string) *models.Doctor
	//
}

type doctorService struct {
	dao *dao.DoctorDao
}

func NewDoctorService() DoctorService {
	return &doctorService{
		dao: dao.NewDoctorDao(dataSource.InstanceMaster()),
	}
}

func (s *doctorService) Add(doctor *models.Doctor) bool {
	ok := s.dao.Add(doctor)

	return ok
}

func (s *doctorService) GetDoctorById(id int) (info *models.Doctor)  {
	return s.dao.GetDoctorById(id)
}

func (s *doctorService) GetDoctorBySpecialty(specialty string) *models.Doctor  {
	return s.dao.GetDoctorBySpecialty(specialty)
}