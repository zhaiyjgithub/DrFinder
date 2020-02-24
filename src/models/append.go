package models

import (
	"time"
)

type Append struct {
	ID        int            `gorm:"column:id;primary_key" json:"AppendID"`
	PostID    int  `gorm:"column:post_id"`
	Content   string `gorm:"column:content"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	CreatedAt time.Time      `gorm:"column:created_at"`
}

// TableName sets the insert table name for this struct type
func (a *Append) TableName() string {
	return "appends"
}
