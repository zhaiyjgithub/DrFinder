package dao

import (
	"DrFinder/src/models/doctor"
	"github.com/jinzhu/gorm"
)

type DoctorDao struct {
	engine *gorm.DB
}

func NewDoctorDao(engine *gorm.DB) *DoctorDao {
	return &DoctorDao{
		engine:engine,
	}
}

func (d *DoctorDao) Add(doctor *doctor.Doctor) error  {
	db := d.engine.Create(doctor)
	return db.Error
}
