package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type AnswerService interface {
	AddAnswer(answer *models.Answer) error
	DeleteByUser(id int, userId int) error
	AddLikes(id int) error
}

type answerService struct {
	dao *dao.AnswerDao
}

func NewAnswerService() AnswerService {
	return &answerService{dao:dao.NewAnswerDao(dataSource.InstanceMaster())}
}

func (s *answerService) AddAnswer(answer *models.Answer) error  {
	return s.AddAnswer(answer)
}

func (s *answerService) DeleteByUser(id int, userId int) error  {
	return s.dao.DeleteByUser(id, userId)
}

func (s *answerService) AddLikes(id int) error   {
	return s.dao.AddLikes(id)
}

