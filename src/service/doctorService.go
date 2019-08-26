package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models/doctor"
)

type DoctorService interface {
	Add(doctor *doctor.Doctor) bool
}

type doctorService struct {
	dao *dao.DoctorDao
}

func NewDoctorService() DoctorService {
	return &doctorService{
		dao: dao.NewDoctorDao(dataSource.InstanceMaster()),
	}
}

func (s *doctorService) Add(doctor *doctor.Doctor) bool {
	ok := s.dao.Add(doctor)

	return ok
}