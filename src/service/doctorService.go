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
	SearchDoctor(doctor *models.Doctor) []models.Doctor
	UpdateDoctorById(doctor *models.Doctor) error
	UpdateDoctorInfo(info *models.Doctor) error
	DeleteDoctorById(id int) bool
	SearchDoctorByPage(doctor *models.Doctor, page int, pageSize int) []models.Doctor
	GetDoctorByPage(page int, pageSize int) []models.Doctor
	GetHotSearchDoctors() *[]models.Doctor
	GetRelatedDoctors(relateDoctor *models.Doctor) *[]models.Doctor
	GetSpecialty() []string
	SearchDoctorsByNpiList(npiList []int) []models.Doctor
	FindDoctorByPage(doctor *models.Doctor, lat float64, lng float64, page int, pageSize int) []models.DoctorGeo
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

func (s *doctorService) SearchDoctor(doctor *models.Doctor) []models.Doctor {
	return s.dao.SearchDoctor(doctor)
}

func (s *doctorService) UpdateDoctorById(doctor *models.Doctor) error  {
	return s.dao.UpdateDoctorById(doctor)
}

func (s *doctorService) UpdateDoctorInfo(info *models.Doctor) error  {
	return s.dao.UpdateDoctorInfo(info)
}

func (s *doctorService) DeleteDoctorById(id int) bool  {
	return s.dao.DeleteDoctorById(id)
}

func (s *doctorService) SearchDoctorByPage(doctor *models.Doctor, page int, pageSize int) []models.Doctor  {
	return s.dao.SearchDoctorByPage(doctor, page, pageSize)
}

func (s *doctorService) GetDoctorByPage(page int, pageSize int) []models.Doctor  {
	return s.dao.GetDoctorByPage(page, pageSize)
}

func (s *doctorService) GetHotSearchDoctors() *[]models.Doctor {
	return s.dao.GetHotSearchDoctors()
}

func (s *doctorService) GetRelatedDoctors(relateDoctor *models.Doctor) *[]models.Doctor {
	return s.dao.GetRelatedDoctors(relateDoctor)
}

func (s *doctorService) GetSpecialty() []string {
	return s.dao.GetSpecialty()
}

func (s *doctorService) SearchDoctorsByNpiList(npiList []int) []models.Doctor  {
	return s.dao.SearchDoctorsByNpiList(npiList)
}

func (s *doctorService) FindDoctorByPage(doctor *models.Doctor, lat float64, lng float64, page int, pageSize int) []models.DoctorGeo   {
	return s.dao.FindDoctorByPage(doctor, lat, lng, page, pageSize)
}