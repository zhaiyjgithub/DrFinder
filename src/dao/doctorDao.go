package dao

import (
	"DrFinder/src/models/doctor"
	"github.com/go-xorm/xorm"
)

type DoctorDao struct {
	engine *xorm.Engine
}

func NewDoctorDao(engine *xorm.Engine) *DoctorDao {
	return &DoctorDao{
		engine:engine,
	}
}

func (d *DoctorDao) Add(doctor *doctor.Doctor) error  {
	_, err := d.engine.Insert(doctor)

	return err
}
