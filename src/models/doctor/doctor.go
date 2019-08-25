package doctor

import "github.com/jinzhu/gorm"

//`doctor_id` int(11) NOT NULL,
//`name_prefix` varchar(5) NULL,
//`first_name` varchar(50) NULL,
//`last_name` varchar(50) NULL,
//`create_date` datetime NOT NULL,
//`npi` int(20) NULL,
//`credential` varchar(10) NULL,
//`gender` varchar(5) NULL,
//`name` varchar(255) NULL,
//`profile` text NULL,
//`middle_name` varchar(20) NULL,
//`mailing_address` varchar(255) NULL,
//`mailing_city` varchar(255) NULL,
//`mailing_state` varchar(10) NULL,
//`mailing_zip_code` varchar(15) NULL,
//`mailing_phone` varchar(15) NULL,
//`mailing_fax` varchar(20) NULL,
//`business_address` varchar(255) NULL,
//`business_city` varchar(255) NULL,
//`business_state` varchar(10) NULL,
//`business_zip_code` varchar(15) NULL,
//`business_phone` varchar(15) NULL,
//`specialty` varchar(255) NULL,

type Doctor struct {
	gorm.Model
	NamePrefix      string
	FirstName       string
	MiddleName      string
	LastName        string
	Name            string
	Npi             int
	Credential      string
	Gender          string
	Specialty       string
	Profile         string
	MailingAddress  string
	MailingCity     string
	MailingState    string
	MailingZipCode  string
	MailingPhone    string
	MailFax         string
	BusinessAddress string
	BusinessCity    string
	BusinessState   string
	BusinessZipCode string
	BusinessPhone   string
}
