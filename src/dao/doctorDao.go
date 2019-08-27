package dao

import (
	"DrFinder/src/models"
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

func (d *DoctorDao) Add(doctor *models.Doctor) bool  {
	ok := d.engine.NewRecord(doctor)
	return ok
}

func (d *DoctorDao) GetDoctorById(id int)  *models.Doctor  {
	var doctor models.Doctor
	db := d.engine.Where("id = ?", id).First(&doctor)

	if db.Error != nil {
		return nil
	}else {
		return &doctor
	}
}

func (d *DoctorDao) GetDoctorBySpecialty(specialty string) *models.Doctor  {
	var doctor models.Doctor
	db:= d.engine.Where("specialty LIKE ?", specialty).First(&doctor)

	if db.Error != nil {
		return  nil
	}

	return &doctor
}
