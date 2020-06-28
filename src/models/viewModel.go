package models

import "time"

type GeoDistance struct {
	Npi int
	Lat float64
	Lng float64
	Distance float64
}

type DoctorGeo struct {
	Lat float64
	Lng float64
	Distance float64
	ID           int
	Npi          int
	LastName     string
	FirstName    string
	MiddleName   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	FullName     string
	NamePrefix   string
	Credential   string
	Gender       string
	Address      string
	AddressSuit  string
	City         string
	State        string
	Zip          string
	Phone        string
	Fax 		 string
	Specialty    string
	SubSpecialty string
	JobTitle     string
	Summary      string
	YearOfExperience string
}
