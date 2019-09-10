package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type PostService interface {
	Add(post *models.Post) error
	GetPostListByPage(post *models.Post, page int, pageSize int) []models.Post
	Delete(id int) error
	Update(newPost *models.Post) error
	AddLike(id int) error
	AddFavor(id int) error
	DeleteByUser(id int, userId int) error
}

type postService struct {
	dao *dao.PostDao

}

func NewPostService() PostService {
	return &postService{dao: dao.NewPostDao(dataSource.InstanceMaster())}
}

func (s *postService) Add(post *models.Post) error  {
	return s.dao.Add(post)
}

func (s *postService) GetPostListByPage(post *models.Post, page int, pageSize int) []models.Post  {
	return s.dao.GetPostListByPage(post, page, pageSize)
}

func (s *postService) Delete(id int) error  {
	return s.dao.Delete(id)
}

func (s *postService) Update(newPost *models.Post) error  {
	return s.dao.Update(newPost)
}

func (s *postService) AddLike(id int) error  {
	return s.dao.AddLike(id)
}

func (s *postService) AddFavor(id int) error  {
	return s.dao.AddFavor(id)
}

func (s *postService) DeleteByUser(id int, userId int) error  {
	return s.dao.DeleteByUser(id, userId)
}