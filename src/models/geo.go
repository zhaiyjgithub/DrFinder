package models

import (
	"time"
)

type Geo struct {
	ID        int             `gorm:"column:id;primary_key"`
	Npi       int   `gorm:"column:npi"`
	Lat       float64 `gorm:"column:lat"`
	Lng       float64 `gorm:"column:lng"`
	CreatedAt time.Time       `gorm:"column:created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at"`
}

// TableName sets the insert table name for this struct type
func (g *Geo) TableName() string {
	return "geos"
}
