package models

import "time"

type UserAction struct {
	Platform string `bson:"platform"`
	Lat float64	`bson:"lat"`
	Lng float64 `bson:"lng"`
	UserID int  `bson:"userId"`
	Name string `bson:"name"`
	AppVersion string `bson:"appVersion"`
	CreatedDate time.Time  `bson:"createdDate"`
}

type UserView struct {
	Platform string `bson:"platform"`
	Lat float64	`bson:"lat"`
	Lng float64 `bson:"lng"`
	UserID int  `bson:"userId"`
	Name string `bson:"name"`
	AppVersion string `bson:"appVersion"`
	BeginTime time.Time `bson:"beginTime"`
	EndTime time.Time	`bson:"endTime"`
}

type UserSearchDrRecord struct {
	Name string `bson:"name"`
	Specialty string `bson:"specialty"`
	Gender string `bson:"gender"`
	City string `bson:"city"`
	State string `bson:"state"`
	Lat float64 `bson:"lat"`
	Lng float64 `bson:"lng"`
	Page int `bson:"page"`
	PageSize int `bson:"pageSize"`
	Platform string `bson:"platform"`
	UserID int `bson:"userId"`
	CreatedDate time.Time `bson:"createDate"`
}

type DrSearchResultRecord struct {
	Npi int `bson:"npi"`
	Lat float64 `bson:"lat"`
	Lng float64 `bson:"lng"`
	Platform string `bson:"platform"`
	UserID int `bson:"userId"`
	CreatedDate time.Time `bson:"createDate"`
}
