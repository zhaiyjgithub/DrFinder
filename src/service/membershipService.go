package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type Membership interface {
	Add(edu *models.Membership) error
	//GetLangByNpi(npi int) *models.Lang
}

type membershipService struct {
	dao *dao.MembershipDao
}

func NewMembershipService() Membership {
	return &membershipService{dao:dao.NewMembershipDao(dataSource.InstanceMaster())}
}

func (s *membershipService) Add(edu *models.Membership) error  {
	return s.dao.Add(edu)
}

func (s *membershipService) GetMemberShipByNpi(npi int) *models.Membership {
	return s.dao.GetMemberShipByNpi(npi)
}







