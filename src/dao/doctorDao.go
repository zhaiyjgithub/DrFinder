package dao

import (
	"DrFinder/src/models/doctor"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type DoctorDao struct {
	engine *gorm.DB
}

func NewDoctorDao(engine *gorm.DB) *DoctorDao {
	return &DoctorDao{
		engine: engine,
	}
}

func (d *DoctorDao) Add(doctor *doctor.Doctor) error {
	info := insert()
	err := d.engine.Create(info)
	fmt.Println(err)

	return nil
}

func insert() *doctor.FavorDoctor {
	info := new(doctor.FavorDoctor)
	info.DoctorId = 34
	info.UserId = 10
	info.CreateAt = time.Now()
	info.UpdateAt = time.Now()

	return info
}
