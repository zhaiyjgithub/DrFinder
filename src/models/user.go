package models

import (
	"time"
)

type User struct {
	ID         int    `gorm:"column:id;primary_key" json:"UserID"`
	Name       string `gorm:"column:name"`
	LastName   string `gorm:"column:last_name"`
	FirstName  string `gorm:"column:first_name"`
	MiddleName string `gorm:"column:middle_name"`
	Email      string `gorm:"column:email"`
	Phone      string `gorm:"column:phone"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"-"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"-"`
	Bio        string `gorm:"column:bio"`
	HeaderIcon string `gorm:"column:header_icon"`
	Level      int  `gorm:"column:level"`
	Password   string `grom:"column:password"`
	Gender     bool `gorm:"column:gender"`
}