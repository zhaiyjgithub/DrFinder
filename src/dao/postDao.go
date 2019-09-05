package dao

import (
	"DrFinder/src/models"
	"github.com/jinzhu/gorm"
)

type PostDao struct {
	engine *gorm.DB
}

func NewPostDao(engine *gorm.DB) *PostDao {
	return &PostDao{engine: engine}
}

func (d *PostDao) Add(post *models.Post) error {
	db := d.engine.Create(post)

	return db.Error
}

func (d *PostDao) GetPostListByPage(post *models.Post, page int, pageSize int) []models.Post {
	var posts []models.Post

	if post.Type > 0 {
		d.engine.Limit(pageSize).Offset((page - 1)*pageSize).Find(&posts, "type = ?", post.Type)
	}else {
		d.engine.Limit(pageSize).Offset((page - 1)*pageSize).Find(&posts)
	}

	return posts
}

func (d *PostDao) Delete(id int) error {
	var post models.Post
	post.ID = id

	db := d.engine.Delete(&post)

	return db.Error
}

func (d *PostDao) Update(newPost *models.Post) error {
	var post models.Post

	db := d.engine.Where("id = ?", newPost.ID).Find(&post)

	if db.Error != nil {
		return db.Error
	}

	db = d.engine.Model(&post).Update(newPost)

	return db.Error
}