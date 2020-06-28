package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type PostImageService interface {
	CreatePostImage(postImage models.PostImage) error
	GetImageByPostId(postId int) []*models.PostImage
}

type postImageService struct {
	dao *dao.PostImageDao
}

func NewPostImageService() PostImageService {
	return &postImageService{dao:dao.NewPostImageDao(dataSource.InstanceMaster())}
}

func (s *postImageService) CreatePostImage(postImage models.PostImage) error {
	return s.dao.CreatePostImage(postImage)
}

func (s *postImageService) GetImageByPostId(postId int) []*models.PostImage {
	return s.dao.GetImageByPostId(postId)
}


