package models

import (
	"time"
)

type Feedback struct {
	ID        int            `gorm:"column:id;primary_key"`
	UserID    int  `gorm:"column:user_id"`
	Content   string `gorm:"column:content"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
}

// TableName sets the insert table name for this struct type
func (f *Feedback) TableName() string {
	return "feedbacks"
}
