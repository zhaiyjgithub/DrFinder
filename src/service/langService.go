package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type LangService interface {
	Add(edu *models.Lang) error
	GetLangByNpi(npi int) *models.Lang
}

type langService struct {
	dao *dao.LangDao
}

func NewLangService() LangService {
	return &langService{dao:dao.NewLangDao(dataSource.InstanceMaster())}
}

func (s *langService) Add(edu *models.Lang) error  {
	return s.dao.Add(edu)
}

func (s *langService) GetLangByNpi(npi int) *models.Lang  {
	return s.dao.GetLangByNpi(npi)
}







