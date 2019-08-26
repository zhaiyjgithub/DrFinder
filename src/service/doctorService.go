package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models/doctorModel"
)

type DoctorService interface {
	Add(doctor *doctorModel.Doctor) bool
	GetDoctorById(id int) *doctorModel.Doctor
}

type doctorService struct {
	dao *dao.DoctorDao
}

func NewDoctorService() DoctorService {
	return &doctorService{
		dao: dao.NewDoctorDao(dataSource.InstanceMaster()),
	}
}

func (s *doctorService) Add(doctor *doctorModel.Doctor) bool {
	ok := s.dao.Add(doctor)

	return ok
}

func (s *doctorService) GetDoctorById(id int) (info *doctorModel.Doctor)  {
	return s.dao.GetDoctorById(id)
}