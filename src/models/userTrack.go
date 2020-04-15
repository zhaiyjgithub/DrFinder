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
