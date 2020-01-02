package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type CollectionService interface {
	Add(userId int, npi int) error
	GetCollections(userId int) []models.Collection
	GetIsHasCollected(userId int, npi int) error
}

type collectionService struct {
	dao *dao.CollectionDao
}

func NewCollectionService() CollectionService  {
	return &collectionService{dao: dao.NewCollectionDao(dataSource.InstanceMaster())}
}

func (s *collectionService) Add(userId int, npi int) error {
	return s.dao.Add(userId, npi)
}

func (s *collectionService) GetCollections(userId int) []models.Collection  {
	return s.dao.GetCollections(userId)
}

func (s *collectionService) GetIsHasCollected(userId int, npi int) error  {
	return s.dao.GetIsHasCollected(userId, npi)
}
