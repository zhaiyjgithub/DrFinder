package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type PostElasticService interface {
	AddOnePost(post *models.Post) error
	QueryPost(content string, page int, pageSize int) []int
}

type postElasticService struct {
	dao *dao.PostElasticDao
}

func NewPostElasticService() PostElasticService {
	return &postElasticService{dao: dao.NewElasticDao(dataSource.InstanceElasticSearchClient())}
}

func (s *postElasticService) AddOnePost(post *models.Post) error {
	return s.dao.AddOnePost(post)
}

func (s *postElasticService) QueryPost(content string, page int, pageSize int) []int {
	return s.dao.QueryPost(content, page, pageSize)
}

