package models

import (
	"time"
)

type Answer struct {
	ID          int            `gorm:"column:id;primary_key"`
	UserID      int  `gorm:"column:user_id"`
	Description string `gorm:"column:description"`
	Likes       int  `gorm:"column:likes"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
}
