package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type CertificationService interface {
	Add(cer *models.Certification) error
	GetCertificationByNpi(npi int) []models.Certification
}

type certificationService struct {
	dao *dao.CertificationDao
}

func NewCertificationService() CertificationService {
	return &certificationService{dao:dao.NewCertificationDao(dataSource.InstanceMaster())}
}

func (s *certificationService) Add(cer *models.Certification) error {
	return s.dao.Add(cer)
}

func (s *certificationService) GetCertificationByNpi(npi int) []models.Certification  {
	return s.dao.GetCertificationByNpi(npi)
}




