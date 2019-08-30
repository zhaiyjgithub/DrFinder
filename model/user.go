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

type User struct {
	ID         int            `gorm:"column:id;primary_key"`
	Name       sql.NullString `gorm:"column:name"`
	Email      sql.NullString `gorm:"column:email"`
	CreatedAt  time.Time      `gorm:"column:created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at"`
	Phone      sql.NullString `gorm:"column:phone"`
	Bio        sql.NullString `gorm:"column:bio"`
	Password   sql.NullString `gorm:"column:password"`
	HeaderIcon sql.NullString `gorm:"column:header_icon"`
	Level      sql.NullInt64  `gorm:"column:level"`
}

// TableName sets the insert table name for this struct type
func (u *User) TableName() string {
	return "users"
}
