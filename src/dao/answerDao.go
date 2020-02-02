package dao

import (
	"DrFinder/src/models"
	"github.com/jinzhu/gorm"
)

type AnswerDao struct {
	engine *gorm.DB
}

func NewAnswerDao(engine *gorm.DB) *AnswerDao {
	return &AnswerDao{engine:engine}
}

func (d *AnswerDao) AddAnswer(answer *models.Answer) error {
	db := d.engine.Create(answer)

	return db.Error
}

func (d *AnswerDao) DeleteByUser(id int, userId int) error {
	db := d.engine.Where("id = ? AND user_id = ?", id, userId).Find(models.Answer{})

	return db.Error
}

func (d *AnswerDao) AddLikes(id int) error {
	var answer models.Answer

	db := d.engine.Where("id = ?", id).Find(&answer)

	if db.Error != nil {
		return db.Error
	}

	answer.Likes = answer.Likes + 1

	db = d.engine.Model(&answer).Update("likes", answer.Likes)

	return db.Error
}

func (d *AnswerDao) GetAnswerListByPage(postId int, page int, pageSize int) []models.Answer {
	var answers []models.Answer

	d.engine.Limit(pageSize).Offset((page - 1)*pageSize).Find(&answers, "post_id = ?", postId)

	return answers
}

func (d *AnswerDao) GetLastAnswer(postId int) *models.Answer {
	var answer models.Answer
	d.engine.Find("post_id = ?", postId).Order("created_at").Last(&answer)

	if answer.PostID == postId {
		return &answer
	} else {
		return nil
	}
}
