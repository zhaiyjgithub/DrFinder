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
		db = d.engine.Create(&models.Collection{UserID: userId, ObjectID: objectID, ObjectType: objectType})
		return db.Error
	}

	return  errors.New("is existing")
}

func (d *CollectionDao) GetCollections(userId int, objectType int) []*models.Collection {
	var collections []*models.Collection

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

func (d *CollectionDao) GetMyFavoriteDoctors(userId int, objectType int, page int, pageSize int) []*models.Doctor {
	var doctors []*models.Doctor
	d.engine.Limit(pageSize).Offset((page - 1)*pageSize).Raw("SELECT * from doctors WHERE npi in " +
		"(SELECT object_id FROM collections WHERE user_id = ? and object_type = ?)", userId, objectType).Find(&doctors)

	return doctors
}

func (d *CollectionDao) GetMyFavoritePosts(userId int, objectType int, page int, pageSize int) []models.Post {
	var posts []models.Post
	d.engine.Limit(pageSize).Offset((page - 1)*pageSize).Raw("select * from posts where id in " +
		"(select object_id from collections where user_id = ? and object_type = ?)", userId, objectType).Find(&posts)

	return posts
}

func (d *CollectionDao) DeleteMyFavorite(userId int, objectIds []int) error {
	db := d.engine.Where("user_id = ? and object_id in (?)", userId, objectIds).Delete(models.Collection{})
	return db.Error
}