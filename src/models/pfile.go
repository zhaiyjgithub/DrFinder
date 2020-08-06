package models

import (
	"database/sql"
	"github.com/guregu/null"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

type Pfile struct {
	ID            int            `gorm:"column:id;primary_key"`
	Npi           int            `gorm:"column:npi"`
	FirstName     string `gorm:"column:first_name"`
	LastName      string `gorm:"column:last_name"`
	MidName       string `gorm:"column:mid_name"`
	FirstAddress  string `gorm:"column:first_address"`
	SecondAddress string `gorm:"column:second_address"`
	City          string `gorm:"column:city"`
	State         string `gorm:"column:state"`
	PostalCode    string `gorm:"column:postal_code"`
	Phone         string `gorm:"column:phone"`
	Fax           string `gorm:"column:fax"`
	Gender        string `gorm:"column:gender"`
}

// TableName sets the insert table name for this struct type
func (p *Pfile) TableName() string {
	return "pfile"
}
