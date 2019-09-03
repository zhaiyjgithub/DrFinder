package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

type Ad struct {
	ID          int            `gorm:"column:id;primary_key"`
	Img         sql.NullString `gorm:"column:img"`
	Title       sql.NullString `gorm:"column:title"`
	Description sql.NullString `gorm:"column:description"`
	Type        sql.NullString `gorm:"column:type"`
	SubID       sql.NullInt64  `gorm:"column:sub_id"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	IsShow      sql.NullInt64  `gorm:"column:is_show"`
	Index       sql.NullInt64  `gorm:"column:index"`
}

// TableName sets the insert table name for this struct type
func (a *Ad) TableName() string {
	return "ad"
}
