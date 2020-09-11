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

func (d *DoctorDao) SearchDoctor(doctor *models.Doctor) []*models.Doctor  {
	var doctors []*models.Doctor
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

func (d *DoctorDao) SearchDoctorByPage(doctor *models.Doctor, page int, pageSize int) []*models.Doctor {
	var doctors []*models.Doctor

	firstName := fmt.Sprintf("%%%s%%", doctor.FirstName)
	lastName := fmt.Sprintf("%%%s%%", doctor.LastName)
	specialty := fmt.Sprintf("%%%s%%", doctor.Specialty)

	groupBy:= "last_name"
	if len(specialty) > 0 {
		groupBy = "specialty"
	}

	d.engine.Limit(pageSize).Offset((page -1)*pageSize).Raw("select * from doctors WHERE id in " +
		"(SELECT id from doctors where " +
		"(last_name like ? or first_name like ?) and specialty like ?" +
		" and gender = ? and city = ? and state = ? group by id) order by ?", lastName, firstName, specialty, doctor.Gender, doctor.City, doctor.State, groupBy).Scan(&doctors)

	return doctors
}

func (d *DoctorDao) FindDoctorByPage(doctor *models.Doctor, lat float64, lng float64, page int, pageSize int) []*models.DoctorGeo  {
	var doctorGeos []*models.DoctorGeo
	var genderList []string

	lastName := fmt.Sprintf("%s%%", doctor.LastName)
	if len(doctor.Gender) == 0 {
		genderList = append(genderList, "F", "M")
	}else {
		genderList = append(genderList, doctor.Gender)
	}

	if len(doctor.LastName) > 0 && len(doctor.Specialty) > 0{//most from filter menu
		eachPageSize := pageSize/2
		d.engine.Raw("select G.lat, G.lng, ACOS(SIN((? * 3.1415) / 180 ) *SIN((lat * 3.1415) / 180 ) +COS((? * 3.1415) / 180 ) * COS((lat * 3.1415) / 180 ) *COS((? * 3.1415) / 180 - (lng * 3.1415) / 180 ) ) * 6380 as distance, DL.*" +
			" from (select * from doctors inner join (select id as Did from doctors where last_name like ? and specialty = ? and state = ? and city = ? and gender in (?) limit ? offset ?) as D " +
			"on D.Did = doctors.id) as DL, geos as G where DL.npi = G.npi",
			lat, lat, lng,
			lastName, doctor.Specialty,
			doctor.State, doctor.City, genderList,
			eachPageSize, (page - 1)*eachPageSize,
			).Scan(&doctorGeos)
	} else if len(doctor.LastName) > 0 && len(doctor.Specialty) == 0 {//most from filter menu
		eachPageSize := pageSize/2
		d.engine.Raw("select G.lat, G.lng, ACOS(SIN((? * 3.1415) / 180 ) *SIN((lat * 3.1415) / 180 ) +COS((? * 3.1415) / 180 ) * COS((lat * 3.1415) / 180 ) *COS((? * 3.1415) / 180 - (lng * 3.1415) / 180 ) ) * 6380 as distance, DL.*" +
			" from (select * from doctors inner join (select id as Did from doctors where last_name like ? and state = ? and city = ? and gender in (?) limit ? offset ?) as D " +
			"on D.Did = doctors.id) as DL, geos as G where DL.npi = G.npi",
			lat, lat, lng,
			lastName,
			doctor.State, doctor.City, genderList,
			eachPageSize, (page - 1)*eachPageSize,
		).Scan(&doctorGeos)
	} else if len(doctor.LastName) == 0 && len(doctor.Specialty) > 0 {// most from specialty menu
		d.engine.Raw("select G.lat, G.lng, ACOS(SIN((? * 3.1415) / 180 ) *SIN((lat * 3.1415) / 180 ) +COS((? * 3.1415) / 180 ) * COS((lat * 3.1415) / 180 ) *COS((? * 3.1415) / 180 - (lng * 3.1415) / 180 ) ) * 6380 as distance, DL.*" +
			" from (select * from doctors inner join (select id as Did from doctors where specialty = ? and state = ? and city = ? and gender in (?) limit ? offset ?) as D " +
			"on D.Did = doctors.id) as DL, geos as G where DL.npi = G.npi",
			lat, lat, lng,
			doctor.Specialty,
			doctor.State, doctor.City, genderList,
			pageSize, (page - 1)*pageSize,
		).Scan(&doctorGeos)
	}else {// near by doctors.
		d.engine.Raw("select G.lat, G.lng, G.distance, D.* from " +
			"(select geos.npi as npi, geos.lat as lat, geos.lng as lng, ACOS(SIN((? * 3.1415) / 180 ) *SIN((lat * 3.1415) / 180 ) + COS((? * 3.1415) / 180 ) * COS((lat * 3.1415) / 180 ) *COS((? * 3.1415) / 180 - (lng * 3.1415) / 180 ) ) * 6380 as distance from geos " +
			"order by distance  LIMIT ? offset ?) G , " +
			"doctors D  where G.npi = D.npi and gender in (?) and city = ? order by distance",
			lat, lat, lng,
			pageSize, (page - 1)*pageSize,
			genderList, doctor.City).Scan(&doctorGeos)
	}

	return doctorGeos
}

func (d *DoctorDao) GetDoctorByPage(page int, pageSize int) []*models.Doctor  {
	var doctors []*models.Doctor
	d.engine.Limit(pageSize).Offset((page - 1) * pageSize).Find(&doctors)
	return doctors
}

func (d *DoctorDao) GetCity() []string {
	var sps []string

	rows, _ := d.engine.Raw("select specialty from doctors where state = 'NY' group by city").Rows()
	defer rows.Close()

	for rows.Next() {
		var name string
		_ = rows.Scan(&name)
		sps = append(sps, name)
	}

	return sps
}

func (d *DoctorDao) GetHotSearchDoctors() []*models.Doctor {
	var doctors []*models.Doctor
	d.engine.Limit(50).Offset(0).Find(&doctors)
	return doctors
}

func (d *DoctorDao) GetRelatedDoctors(relateDoctor *models.Doctor) []*models.Doctor {
	var doctors []*models.Doctor
	d.engine.Raw("select * from doctors where specialty = ? and npi != ? and city = ? and state = ?",
		relateDoctor.Specialty,
		relateDoctor.Npi,
		relateDoctor.City,
		relateDoctor.State,
		).Scan(&doctors)
	return doctors
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
		_ = rows.Scan(&name)
		sps = append(sps, name)
	}

	return sps
}

func (d *DoctorDao) SearchDoctorsByNpiList(npiList []int) []*models.Doctor  {
	var doctors []*models.Doctor
	d.engine.Raw("select * from doctors where npi in (?)", npiList).Scan(&doctors)
	return doctors
}

func (d *DoctorDao) GetDoctorByNpi(npi int) *models.Doctor {
	var doctor *models.Doctor
	d.engine.Raw("select * from doctors where npi = ?", npi).Scan(&doctor)
	return doctor
}

func (d *DoctorDao) GetDoctorByNpiList(npiList []int) []*models.Doctor {
	doctors := make([]*models.Doctor, 0)
	d.engine.Raw("select * from doctors where npi in (?)", npiList).Find(&doctors)

	return doctors
}

func (d *DoctorDao) GetDoctorsNoAddress(page int , pageSize int) []*models.Doctor  {
	var docs []*models.Doctor
	d.engine.Raw("select * from doctors where address = '' limit ? offset ?", pageSize, (page - 1)*pageSize).Scan(&docs)
	return docs
}

func (d *DoctorDao) UpdateDoctorAddress(doc *models.Doctor) error  {
	db := d.engine.Model(&doc).Update("address", doc.Address)
	return db.Error
}