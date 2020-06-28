package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type AwardService interface {
	Add(award *models.Award) error
	GetAwardByNpi(npi int) []*models.Award
}

type awardService struct {
	dao *dao.AwardDao
}

func NewAwardService() AwardService {
	return &awardService{
		dao:dao.NewAwardDao(dataSource.InstanceMaster()),
	}
}

func (s *awardService) Add(award *models.Award) error  {
	return s.dao.Add(award)
}

func (s *awardService) GetAwardByNpi(npi int) []*models.Award  {
	return s.dao.GetAwardByNpi(npi)
}
