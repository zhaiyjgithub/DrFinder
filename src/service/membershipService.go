package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type MembershipService interface {
	Add(edu *models.Membership) error
	GetMemberShipByNpi(npi int) []*models.Membership
}

type membershipService struct {
	dao *dao.MembershipDao
}

func NewMembershipService() MembershipService {
	return &membershipService{dao:dao.NewMembershipDao(dataSource.InstanceMaster())}
}

func (s *membershipService) Add(edu *models.Membership) error  {
	return s.dao.Add(edu)
}

func (s *membershipService) GetMemberShipByNpi(npi int) []*models.Membership {
	return s.dao.GetMemberShipByNpi(npi)
}







