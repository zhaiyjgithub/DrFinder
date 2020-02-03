package models

type PostImage struct {
	ID     int  `gorm:"column:id;primary_key" json:"-"`
	PostID int  `gorm:"column:post_id" json:"-"`
	URL    string `gorm:"column:url"`
}

// TableName sets the insert table name for this struct type
func (p *PostImage) TableName() string {
	return "post_images"
}
