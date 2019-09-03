package dao

import (
	"DrFinder/src/models"
	"github.com/jinzhu/gorm"
)

type AdvertisementDao struct {
	engine *gorm.DB
}

func NewAdvertisementDao(engine *gorm.DB) *AdvertisementDao {
	return &AdvertisementDao{engine: engine}
}

func (d *AdvertisementDao) Add(ad *models.Advertisement) error {
	db := d.engine.Create(ad)

	return db.Error
}

func (d *AdvertisementDao) GetAdvertisementList() *[]models.Advertisement {
	var ads []models.Advertisement
	d.engine.Find(&ads)

	return &ads
}

func (d *AdvertisementDao) Update(newAd *models.Advertisement) error {
	var ad models.Advertisement

	db := d.engine.Where("id = ?", newAd.ID).First(&ad)

	if db.Error != nil {
		return db.Error
	}

	db = d.engine.Model(&ad).Update(newAd)

	if db.Error != nil {
		return db.Error
	}

	return nil
}

func (d *AdvertisementDao) Delete(id int) error {
	var ad models.Advertisement
	ad.ID = id

	db := d.engine.Delete(&ad)

	return db.Error
}
