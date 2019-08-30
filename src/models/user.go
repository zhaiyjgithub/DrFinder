package models

import (
	"time"
)

type User struct {
	ID         int            `gorm:"column:id;primary_key"`
	Name       string `gorm:"column:name"`
	Email      string `gorm:"column:email"`
	CreatedAt  time.Time      `gorm:"column:created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at"`
	Phone      string `gorm:"column:phone"`
	Bio        string `gorm:"column:bio"`
	Password   string `gorm:"column:password"`
	HeaderIcon string `gorm:"column:header_icon"`
	Level      int  `gorm:"column:level"`
}