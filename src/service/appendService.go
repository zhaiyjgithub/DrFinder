package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type AppendService interface {
	AddAppend(append *models.Append) error
	GetAppends(postId int) []models.Append
}

type appendService struct {
	dao *dao.AppendDao
}

func NewAppendService() AppendService {
	return &appendService{dao:dao.NewAppendDao(dataSource.InstanceMaster())}
}

func (s *appendService)AddAppend(append *models.Append) error  {
	return s.dao.AddAppend(append)
}

func (s *appendService)GetAppends(postId int) []models.Append  {
	return s.dao.GetAppends(postId)
}


