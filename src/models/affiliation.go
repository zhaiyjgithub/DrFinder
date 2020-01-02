package models

import (
	"time"
)

type Affiliation struct {
	ID        int            `gorm:"column:id;primary_key" json:"-"`
	Npi       int `gorm:"column:npi" json:"-"`
	Name      string `gorm:"column:name"`
	Desc      string `gorm:"column:desc"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"-"`
}

// TableName sets the insert table name for this struct type
func (a *Affiliation) TableName() string {
	return "affiliations"
}
