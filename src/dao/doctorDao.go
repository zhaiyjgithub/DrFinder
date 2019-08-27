package dao

import (
	"DrFinder/src/models/doctorModel"
	"github.com/jinzhu/gorm"
)

type DoctorDao struct {
	engine *gorm.DB
}

func NewDoctorDao(engine *gorm.DB) *DoctorDao {
	return &DoctorDao{
		engine: engine,
	}
}

func (d *DoctorDao) Add(doctor *doctorModel.Doctor) bool  {
	ok := d.engine.NewRecord(doctor)
	return ok
}

func (d *DoctorDao) GetDoctorById(id int) *doctorModel.Doctor  {
	 var info = new(doctorModel.Doctor)
	 db := d.engine.Where("id =  ?", id).First(info)

	 if db.Error != nil {
	 	return nil
	 }else {
	 	return info
	 }
}
