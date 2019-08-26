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

type Doctor struct {
	ID              int            `gorm:"column:id;primary_key"`
	Npi             sql.NullInt64  `gorm:"column:npi"`
	LastName        sql.NullString `gorm:"column:last_name"`
	FirstName       sql.NullString `gorm:"column:first_name"`
	MiddleName      sql.NullString `gorm:"column:middle_name"`
	CreatedAt       time.Time      `gorm:"column:created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at"`
	Name            sql.NullString `gorm:"column:name"`
	NamePrefix      sql.NullString `gorm:"column:name_prefix"`
	Credential      sql.NullString `gorm:"column:credential"`
	Gender          sql.NullString `gorm:"column:gender"`
	MailingAddress  sql.NullString `gorm:"column:mailing_address"`
	MailingCity     sql.NullString `gorm:"column:mailing_city"`
	MailingState    sql.NullString `gorm:"column:mailing_state"`
	MailingZip      sql.NullString `gorm:"column:mailing_zip"`
	MailingPhone    sql.NullString `gorm:"column:mailing_phone"`
	MailingFax      sql.NullString `gorm:"column:mailing_fax"`
	BusinessAddress sql.NullString `gorm:"column:business_address"`
	BusinessCity    sql.NullString `gorm:"column:business_city"`
	BusinessState   sql.NullString `gorm:"column:business_state"`
	BusinessZip     sql.NullString `gorm:"column:business_zip"`
	BusinessPhone   sql.NullString `gorm:"column:business_phone"`
	BusinessFax     sql.NullString `gorm:"column:business_fax"`
	Specialty       sql.NullString `gorm:"column:specialty"`
}

// TableName sets the insert table name for this struct type
func (d *Doctor) TableName() string {
	return "doctors"
}
