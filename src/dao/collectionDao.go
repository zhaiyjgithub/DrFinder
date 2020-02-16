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

func (d *CollectionDao) Add(userId int, objectID int, objectType int) error {
	var collection models.Collection

	db := d.engine.Where(&models.Collection{UserID: userId, ObjectID: objectID, ObjectType: objectType}).Find(&collection)
	if db.Error != nil {
		db = d.engine.Create(&models.Collection{UserID: userId, ObjectID: objectID})
		return db.Error
	}

	return  errors.New("is existing")
}

func (d *CollectionDao) GetCollections(userId int, objectType int) []models.Collection {
	var collections []models.Collection

	db := d.engine.Where("user_id = ? AND object_type = ?", userId, objectType).Find(&collections)

	if db.Error != nil {
		return nil
	}

	return collections
}

func (d *CollectionDao) GetIsHasCollected(userId int, objectID int, objectType int) error {
	var collection models.Collection

	db := d.engine.Where("user_id = ? AND object_id = ? AND object_type = ?", userId, objectID, objectType).Find(&collection)

	return db.Error
}

func (d *CollectionDao) Delete(userId int, objectID int, objectType int) error {
	var collection models.Collection

	db := d.engine.Where(&models.Collection{UserID: userId, ObjectID: objectID, ObjectType: objectType}).Find(&collection)

	if db.Error != nil {
		return db.Error
	}

	db = db.Delete(&collection)

	return db.Error
}
