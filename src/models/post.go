package models

import (
	"time"
)

type Post struct {
	ID          int            `gorm:"column:id;primary_key"`
	Type        int  `gorm:"column:type"`
	Priority    int  `gorm:"column:priority"`
	Title       string `gorm:"column:title"`
	UserID      int `gorm:"column:user_id"`
	Description string `gorm:"column:description"`
	Favorites   int  `gorm:"column:favorites"`
	Likes       int  `gorm:"column:likes"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"-"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"-"`
}
