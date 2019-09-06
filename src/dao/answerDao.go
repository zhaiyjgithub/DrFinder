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
