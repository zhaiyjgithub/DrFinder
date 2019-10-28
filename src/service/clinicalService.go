package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type ClinicalService interface {
	Add(clinic *models.Clinical) error
	GetClinicalByNpi(npi int) *models.Clinical
}

type clinicService struct {
	dao *dao.ClinicalDao
}

func NewClinicalService() ClinicalService{
	return &clinicService{dao:dao.NewClinicalDao(dataSource.InstanceMaster())}
}

func (s *clinicService) Add(clinic *models.Clinical) error  {
	return s.dao.Add(clinic)
}

func (s *clinicService) GetClinicalByNpi(npi int) *models.Clinical  {
	return s.dao.GetClinicalByNpi(npi)
}
