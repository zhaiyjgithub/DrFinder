package models

import "time"

type Advertisement struct {
	ID          int            `gorm:"column:id;primary_key"`
	Img         string 			`gorm:"column:img"`
	Title       string 			`gorm:"column:title"`
	Description string 			`gorm:"column:description"`
	Type        string 			`gorm:"column:type"`
	SubID       int  			`gorm:"column:sub_id"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	IsShow      int  			`gorm:"column:is_show"`
	Index       int 			`gorm:"column:index"`
}
