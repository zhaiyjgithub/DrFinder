package models

import (
	"time"
)

type Doctor struct {
	ID           int            `gorm:"column:id;primary_key"`
	Npi          int  `gorm:"column:npi"`
	LastName     string `gorm:"column:last_name"`
	FirstName    string `gorm:"column:first_name"`
	MiddleName   string `gorm:"column:middle_name"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	FullName     string `gorm:"column:full_name"`
	NamePrefix   string `gorm:"column:name_prefix"`
	Credential   string `gorm:"column:credential"`
	Gender       string `gorm:"column:gender"`
	Address      string `gorm:"column:address"`
	City         string `gorm:"column:city"`
	State        string `gorm:"column:state"`
	Zip          string `gorm:"column:zip"`
	Phone        string `gorm:"column:phone"`
	Fax 		 string `gorm:"column:fax"`
	Specialty    string `gorm:"column:specialty"`
	SubSpecialty string `gorm:"column:sub_specialty"`
	JobTitle     string `gorm:"column:job_title"`
	Summary      string `gorm:"column:summary"`
}

// TableName sets the insert table name for this struct type
func (d *Doctor) TableName() string {
	return "doctors"
}
