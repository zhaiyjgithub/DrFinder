package dao

import (
	"DrFinder/src/models"
	"github.com/jinzhu/gorm"
)

type PostImageDao struct {
	engine *gorm.DB
}

func NewPostImageDao(engine *gorm.DB) *PostImageDao {
	return &PostImageDao{engine:engine}
}

func (d *PostImageDao) CreatePostImage(postImage models.PostImage) error {
	db := d.engine.Create(&postImage)

	return db.Error
}

func (d *PostImageDao) GetImageByPostId(postId int) []*models.PostImage {
	var postImages []*models.PostImage

	d.engine.Where("post_id = ?", postId).Find(&postImages)

	return postImages
}


