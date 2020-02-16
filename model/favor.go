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

type Favor struct {
	ID       int           `gorm:"column:id;primary_key"`
	UserID   sql.NullInt64 `gorm:"column:user_id"`
	Type     sql.NullInt64 `gorm:"column:type"`
	ObjectID sql.NullInt64 `gorm:"column:object_id"`
}

// TableName sets the insert table name for this struct type
func (f *Favor) TableName() string {
	return "favors"
}
