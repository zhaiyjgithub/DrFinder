package dao

import (
	"DrFinder/src/models"
	"fmt"
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

func (d *DoctorDao) SearchDoctorByPage(doctor *models.Doctor, page int, pageSize int) []models.Doctor {
	var doctors []models.Doctor

	firstName := fmt.Sprintf("%%%s%%", doctor.FirstName)
	lastName := fmt.Sprintf("%%%s%%", doctor.LastName)
	specialty := fmt.Sprintf("%%%s%%", doctor.Specialty)
	gender := fmt.Sprintf("%%%s%%", doctor.Gender)
	businessCity := fmt.Sprintf("%%%s%%", doctor.City)

	d.engine.Limit(pageSize).Offset((page -1)*pageSize).Find(&doctors, "first_Name LIKE ? " +
		"AND last_name LIKE ? " +
		"AND specialty LIKE ? " +
		"AND gender LIKE ? " +
		"AND business_city LIKE ?", firstName, lastName, specialty, gender, businessCity)

	return doctors
}

func (d *DoctorDao) GetDoctorByPage(page int, pageSize int) []models.Doctor  {
	var doctors []models.Doctor

	d.engine.Limit(pageSize).Offset((page - 1) * pageSize).Find(&doctors)

	return doctors
}
