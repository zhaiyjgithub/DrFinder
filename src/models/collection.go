package models

import (
	"time"
)

type Collection struct {
	ID        int           `gorm:"column:id;primary_key" json:"-"`
	Npi       int `gorm:"column:npi"`
	UserID    int `gorm:"column:user_id"`
	CreatedAt time.Time     `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time     `gorm:"column:updated_at" json:"-"`
}

// TableName sets the insert table name for this struct type
func (c *Collection) TableName() string {
	return "collections"
}
