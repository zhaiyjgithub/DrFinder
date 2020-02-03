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
	GetAnswerListByPage(postId int, page int, pageSize int) []models.Answer
	GetLastAnswer(postId int) (*models.Answer, int)
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

func (s *answerService) GetAnswerListByPage(postId int, page int, pageSize int) []models.Answer  {
	return s.dao.GetAnswerListByPage(postId, page, pageSize)
}

func (s *answerService) GetLastAnswer(postId int) (*models.Answer, int)  {
	return s.dao.GetLastAnswer(postId)
}
