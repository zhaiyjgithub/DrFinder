package dao

import (
	"DrFinder/src/models"
	"errors"
	"github.com/jinzhu/gorm"
)

type CollectionDao struct {
	engine *gorm.DB
}

func NewCollectionDao(engine *gorm.DB) *CollectionDao {
	return &CollectionDao{engine:engine}
}

func (d *CollectionDao) Add(userId int, npi int) error {
	var collection models.Collection

	db := d.engine.Where(&models.Collection{UserID: userId, Npi: npi}).Find(&collection)
	if db.Error != nil {
		db = d.engine.Create(&models.Collection{UserID: userId, Npi: npi})
		return db.Error
	}

	return  errors.New("is existing")
}

func (d *CollectionDao) GetCollections(userId int) []models.Collection {
	var collections []models.Collection

	db := d.engine.Where("user_id = ?", userId).Find(&collections)

	if db.Error != nil {
		return nil
	}

	return collections
}

func (d *CollectionDao) GetIsHasCollected(userId int, npi int) error {
	var collection models.Collection

	db := d.engine.Where("user_id = ? AND npi = ?", userId, npi).Find(&collection)

	return db.Error
}

func (d *CollectionDao) Delete(userId int, npi int) error {
	var collection models.Collection

	db := d.engine.Where(&models.Collection{UserID: userId, Npi: npi}).Find(&collection)

	if db.Error != nil {
		return db.Error
	}

	db = db.Delete(&collection)

	return db.Error
}
