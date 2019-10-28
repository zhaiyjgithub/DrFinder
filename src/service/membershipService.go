package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type Membership interface {
	Add(edu *models.Lang) error
	GetLangByNpi(npi int) *models.Lang
}

type membership struct {
	dao *dao.LangDao
}

func NewMembership() Membership {
	return &membership{dao:dao.NewLangDao(dataSource.InstanceMaster())}
}

func (s *membership) Add(edu *models.Lang) error  {
	return s.dao.Add(edu)
}

func (s *membership) GetLangByNpi(npi int) *models.Lang  {
	return s.dao.GetLangByNpi(npi)
}







