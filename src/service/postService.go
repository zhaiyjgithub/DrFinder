package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type PostService interface {
	Add(post *models.Post) (error, int)
	GetPostListByPage(postType int, page int, pageSize int) []*models.Post
	Delete(id int) error
	Update(newPost *models.Post) error
	AddLike(id int) error
	AddFavor(id int) error
	DeleteByUser(id int, userId int) error
	GetMyPostListByPage(userId int, page int, pageSize int) []*models.Post
}


type postService struct {
	dao *dao.PostDao
}

func NewPostService() PostService {
	return &postService{dao: dao.NewPostDao(dataSource.InstanceMaster())}
}

func (s *postService) Add(post *models.Post) (error, int)  {
	return s.dao.Add(post)
}

func (s *postService) GetPostListByPage(postType int, page int, pageSize int) []*models.Post  {
	return s.dao.GetPostListByPage(postType, page, pageSize)
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

func (s *postService) GetMyPostListByPage(userId int, page int, pageSize int) []*models.Post {
	return s.dao.GetMyPostListByPage(userId, page, pageSize)
}