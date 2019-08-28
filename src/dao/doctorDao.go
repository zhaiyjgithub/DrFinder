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
	db := d.engine.Create(doctor)
	if db.Error != nil {
		return false
	}
	return true
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

func (d *DoctorDao) SearchDoctor(doctor *models.Doctor) []models.Doctor  {
	var doctors []models.Doctor
	d.engine.Where(doctor).Find(&doctors)

	return doctors
}

func (d *DoctorDao) UpdateDoctorById(newDoctor *models.Doctor) error {
	var doctor models.Doctor
	db:= d.engine.Where("id = ?", newDoctor.ID).First(&doctor)

	if db.Error != nil {
		return db.Error
	}

	db = d.engine.Model(&doctor).Update(newDoctor)

	return db.Error
}

func (d *DoctorDao) DeleteDoctorById(id int) bool {
	var doctor models.Doctor
	doctor.ID = id

	if doctor.ID > 0 {
		db:= d.engine.Delete(&doctor)

		if db.Error != nil {
			return false
		}

		return true
	}

	return false
}
