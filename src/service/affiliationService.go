package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type AffiliationService interface {
	Add(affiliation *models.Affiliation) error
	GetAffiliationByNpi(npi int) []models.Affiliation
	GetAll(page int, pageSize int) error
}

type affiliationService struct {
	dao *dao.AffiliationDao
}

func NewAffiliationService() AffiliationService {
	return &affiliationService{
		dao:dao.NewAffiliationDao(dataSource.InstanceMaster()),
	}
}

func (s *affiliationService)Add(affiliation *models.Affiliation) error  {
	return s.dao.Add(affiliation)
}

func (s *affiliationService) GetAffiliationByNpi(npi int) []models.Affiliation  {
	return s.dao.GetAffiliationByNpi(npi)
}

func (s *affiliationService) GetAll(page int, pageSize int) error  {
	return s.dao.GetAll(page, pageSize)
}
