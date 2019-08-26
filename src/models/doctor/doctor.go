package doctor

import "time"

type Doctor struct {
	ID              int            `gorm:"column:id;primary_key"`
	Npi             int64  `gorm:"column:npi"`
	LastName        string `gorm:"column:last_name"`
	FirstName       string `gorm:"column:first_name"`
	MiddleName      string `gorm:"column:middle_name"`
	CreatedAt       time.Time      `gorm:"column:created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at"`
	Name            string `gorm:"column:name"`
	NamePrefix      string `gorm:"column:name_prefix"`
	Credential      string `gorm:"column:credential"`
	Gender          string `gorm:"column:gender"`
	MailingAddress  string `gorm:"column:mailing_address"`
	MailingCity     string `gorm:"column:mailing_city"`
	MailingState    string `gorm:"column:mailing_state"`
	MailingZip      string `gorm:"column:mailing_zip"`
	MailingPhone    string `gorm:"column:mailing_phone"`
	MailingFax      string `gorm:"column:mailing_fax"`
	BusinessAddress string `gorm:"column:business_address"`
	BusinessCity    string `gorm:"column:business_city"`
	BusinessState   string `gorm:"column:business_state"`
	BusinessZip     string `gorm:"column:business_zip"`
	BusinessPhone   string `gorm:"column:business_phone"`
	BusinessFax     string `gorm:"column:business_fax"`
	Specialty       string `gorm:"column:specialty"`
}