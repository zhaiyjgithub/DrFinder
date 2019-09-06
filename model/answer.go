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

type Answer struct {
	ID          int            `gorm:"column:id;primary_key"`
	UserID      sql.NullInt64  `gorm:"column:user_id"`
	Description sql.NullString `gorm:"column:description"`
	Likes       sql.NullInt64  `gorm:"column:likes"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
}

// TableName sets the insert table name for this struct type
func (a *Answer) TableName() string {
	return "answers"
}
