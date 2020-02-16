package models

type Favor struct {
	ID       int           `gorm:"column:id;primary_key"`
	UserID   int `gorm:"column:user_id"`
	Type     int `gorm:"column:type"`
	ObjectID int `gorm:"column:object_id"`
}

// TableName sets the insert table name for this struct type
func (f *Favor) TableName() string {
	return "favors"
}
