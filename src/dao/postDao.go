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

func (d *PostDao) Add(post *models.Post) (error, int) {
	db := d.engine.Create(post)

	return db.Error, post.ID
}

func (d *PostDao) GetPostListByPage(postType int, page int, pageSize int) []*models.Post {
	var posts []*models.Post

	if postType > 0 {
		d.engine.Limit(pageSize).Offset((page - 1)*pageSize).Order("created_at desc").Find(&posts, "type = ?", postType)
	}else {
		d.engine.Limit(pageSize).Offset((page - 1)*pageSize).Order("created_at desc").Find(&posts)
	}

	return posts
}

func (d *PostDao) GetPostByPostId(ids []int) []*models.Post {
	var post []*models.Post
	d.engine.Where("id in (?)", ids).Find(&post)

	return post
}

func (d *PostDao) GetMyPostListByPage(userId int, page int, pageSize int) []*models.Post {
	var posts []*models.Post
	d.engine.Limit(pageSize).Offset((page - 1)*pageSize).Where("user_id = ?", userId).Order("title", false).Find(&posts)

	return posts
}

func (d *PostDao) Delete(id int) error {
	var post models.Post
	post.ID = id

	db := d.engine.Delete(&post)

	return db.Error
}

func (d *PostDao) DeleteByUser(id int, userId int) error {
	db := d.engine.Where("id = ? AND user_id = ?", id, userId).Delete(models.Post{})

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

func (d *PostDao) AddLike(id int) error {
	var post models.Post

	db := d.engine.Where("id = ?", id).Find(&post)

	if db.Error != nil {
		return db.Error
	}

	post.Likes = post.Likes + 1
	db = d.engine.Model(&post).Update("likes", post.Likes)

	return db.Error
}

func (d *PostDao) AddFavor(id int) error {
	var post models.Post

	db := d.engine.Where("id = ?", id).Find(&post)

	if db.Error != nil {
		return db.Error
	}

	post.Favorites = post.Favorites + 1
	db = d.engine.Model(&post).Update("favorites", post.Favorites)

	return db.Error
}