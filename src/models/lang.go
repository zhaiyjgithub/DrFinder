package models

import (
	"time"
)

type Lang struct {
	ID        int            `gorm:"column:id;primary_key" json:"-"`
	Npi       int  `gorm:"column:npi" json:"-"`
	Lang      string `gorm:"column:lang"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"-"`
}

// TableName sets the insert table name for this struct type
func (l *Lang) TableName() string {
	return "langs"
}
