package dao

import (
	"DrFinder/src/models"
	"github.com/jinzhu/gorm"
)

type FeedbackDao struct {
	engine *gorm.DB
}

func NewFeedbackDao(engine *gorm.DB) *FeedbackDao {
	return &FeedbackDao{engine: engine}
}

func (d *FeedbackDao) AddFeedback(feedback *models.Feedback) error {
	db := d.engine.Create(feedback)

	return db.Error
}