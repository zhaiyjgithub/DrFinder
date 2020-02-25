package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type FeedbackService interface {
	AddFeedback(feedback *models.Feedback) error
}

type feedbackService struct {
	dao *dao.FeedbackDao
}

func NewFeedbackService() FeedbackService {
	return &feedbackService{dao:dao.NewFeedbackDao(dataSource.InstanceMaster())}
}

func (s *feedbackService)AddFeedback(feedback *models.Feedback) error  {
	return s.dao.AddFeedback(feedback)
}
