package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type UserTrackService interface {
	AddActionEvent(actEvent *models.UserAction) error
	FindActionEvent(filter interface{}) []models.UserAction
}

type userTrackService struct {
	dao *dao.UserTrackDao
}

func NewUserTrackService() UserTrackService {
	return &userTrackService{dao:dao.NewUserTrackDao(dataSource.InstanceMongoDB())}
}

func (s *userTrackService) AddActionEvent(actEvent *models.UserAction) error  {
	return s.dao.AddActionEvent(actEvent)
}

func (s *userTrackService) FindActionEvent(filter interface{}) []models.UserAction  {
	return s.dao.FindActionEvent(filter)
}