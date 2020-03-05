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

func (d *DoctorDao) UpdateDoctorInfo(info *models.Doctor) error  {
	var doctor models.Doctor
	db := d.engine.Where("id = ?", info.ID).First(&doctor)

	if db.Error != nil {
		return db.Error
	}

	db = d.engine.Model(&doctor).Update(info)

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
	city := fmt.Sprintf("%%%s%%", doctor.City)
	state := fmt.Sprintf("%%%s%%", doctor.State)

	groupBy:= "last_name"
	if len(specialty) > 0 {
		groupBy = "specialty"
	}

	//$sql='select * from users_location where latitude > '.$lat.'-1 and latitude < '.$lat.'+1 and longitude > '.$lon.'-1 and longitude < '.$lon.'+1 order by ACOS(SIN(('.$lat.' * 3.1415) / 180 ) *SIN((latitude * 3.1415) / 180 ) +COS(('.$lat.' * 3.1415) / 180 ) * COS((latitude * 3.1415) / 180 ) *COS(('.$lon.'* 3.1415) / 180 - (longitude * 3.1415) / 180 ) ) * 6380 asc limit 10';
	//select npi, lat, lng, ACOS(SIN((33.506493 * 3.1415) / 180 ) *SIN((lat * 3.1415) / 180 ) +COS((33.506493 * 3.1415) / 180 ) * COS((lat * 3.1415) / 180 ) *COS((-86.77556* 3.1415) / 180 - (lng * 3.1415) / 180 ) ) * 6380 as distance  from geos where lat > (33.506493 - 1) and lat < (33.506493 + 1) and lng > (-86.77556 - 1) and lng < (-86.77556 + 1)  order by distance asc LIMIT 100 OFFSET 0
	//使用高级联结查询 UNION
	d.engine.Limit(pageSize).Offset((page -1)*pageSize).Raw("select * from doctors WHERE id in " +
		"(SELECT id from doctors where " +
		"(last_name like ? or first_name like ?) and specialty like ?" +
		" and gender like ? and city like ? and state like ? group by id) order by ?", lastName, firstName, specialty, gender, city, state, groupBy).Scan(&doctors)

	return doctors
}

func (d *DoctorDao) GetDoctorByPage(page int, pageSize int) []models.Doctor  {
	var doctors []models.Doctor

	d.engine.Limit(pageSize).Offset((page - 1) * pageSize).Find(&doctors)

	return doctors
}

func (d *DoctorDao) GetHotSearchDoctors() *[]models.Doctor {
	var doctors []models.Doctor

	d.engine.Limit(50).Offset(0).Find(&doctors)

	return &doctors
}

func (d *DoctorDao) GetRelatedDoctors(relateDoctor *models.Doctor) *[]models.Doctor {
	var doctors []models.Doctor

	d.engine.Limit(10).Offset(100).Find(&doctors)
	return &doctors
}

func (d *DoctorDao) GetDoctorStarStatus(userId int, npi int) bool {
	return false
}

func (d *DoctorDao) GetSpecialty() []string {
	var sps []string

	rows, _ := d.engine.Raw("select specialty from doctors").Rows()
	defer rows.Close()

	for rows.Next() {
		var name string
		rows.Scan(&name)

		sps = append(sps, name)
	}

	return sps
}