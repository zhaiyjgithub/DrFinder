package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type EducationService interface {
	Add(edu *models.Education) error
	GetEducationByNpi(npi int) []*models.Education
}

type educationService struct {
	dao *dao.EducationDao
}

func NewEducationService() EducationService {
	return &educationService{dao:dao.NewEducationDao(dataSource.InstanceMaster())}
}

func (s *educationService) Add(edu *models.Education) error  {
	return s.dao.Add(edu)
}

func (s *educationService) GetEducationByNpi(npi int) []*models.Education  {
	return s.dao.GetEducationByNpi(npi)
}






