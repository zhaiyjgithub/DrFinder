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

type Post struct {
	ID          int            `gorm:"column:id;primary_key"`
	Type        sql.NullInt64  `gorm:"column:type"`
	Priority    sql.NullInt64  `gorm:"column:priority"`
	Title       sql.NullString `gorm:"column:title"`
	UsrID       int            `gorm:"column:usr_id"`
	Description sql.NullString `gorm:"column:description"`
	Favorites   sql.NullInt64  `gorm:"column:favorites"`
	Likes       sql.NullInt64  `gorm:"column:likes"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
}

// TableName sets the insert table name for this struct type
func (p *Post) TableName() string {
	return "posts"
}
