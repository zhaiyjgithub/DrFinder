package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type UserTrackService interface {
	AddActionEvent(event *models.UserAction) error
	FindActionEvent(filter interface{}) []models.UserAction
	AddManyActionEvent(events []models.UserAction) error
	AddViewEvent(event *models.UserView) error
	AddManyViewTimeEvent(events []models.UserView) error
	AddSearchDrsRecord(record *models.UserSearchDrRecord) error
}

type userTrackService struct {
	dao *dao.UserTrackDao
}

func NewUserTrackService() UserTrackService {
	return &userTrackService{dao:dao.NewUserTrackDao(dataSource.InstanceMongoDB())}
}

func (s *userTrackService) AddActionEvent(event *models.UserAction) error  {
	return s.dao.AddActionEvent(event)
}

func (s *userTrackService) AddManyActionEvent(events []models.UserAction) error {
	return s.dao.AddManyActionEvent(events)
}

func (s *userTrackService) AddViewEvent(event *models.UserView) error {
	return s.dao.AddViewEvent(event)
}

func (s *userTrackService) AddManyViewTimeEvent(events []models.UserView) error {
	return s.dao.AddManyViewTimeEvent(events)
}

func (s *userTrackService) AddSearchDrsRecord(record *models.UserSearchDrRecord) error  {
	return s.dao.AddSearchDrsRecord(record)
}

func (s *userTrackService) FindActionEvent(filter interface{}) []models.UserAction  {
	return s.dao.FindActionEvent(filter)
}