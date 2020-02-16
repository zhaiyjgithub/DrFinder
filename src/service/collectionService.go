package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type CollectionService interface {
	Add(userId int, objectID int, objectType int) error
	GetCollections(userId int, objectType int) []models.Collection
	GetIsHasCollected(userId int, objectID int, objectType int) error
	Delete(userId int, objectID int, objectType int) error
}

type collectionService struct {
	dao *dao.CollectionDao
}

func NewCollectionService() CollectionService  {
	return &collectionService{dao: dao.NewCollectionDao(dataSource.InstanceMaster())}
}

func (s *collectionService) Add(userId int, objectID int, objectType int) error {
	return s.dao.Add(userId, objectID, objectType)
}

func (s *collectionService) GetCollections(userId int, objectType int) []models.Collection  {
	return s.dao.GetCollections(userId, objectType)
}

func (s *collectionService) GetIsHasCollected(userId int, objectID int, objectType int) error  {
	return s.dao.GetIsHasCollected(userId, objectID, objectType)
}

func (s *collectionService) Delete(userId int, objectID int, objectType int) error  {
	return s.dao.Delete(userId, objectID, objectType)
}
