package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type AdvertisementService interface {
	Add(ad *models.Advertisement) error
	GetAdvertisementList() *[]models.Advertisement
	Update(newAd *models.Advertisement) error
	Delete(id int) error
}

type advertisementService struct {
	dao *dao.AdvertisementDao
}

func NewAdvertiseService() AdvertisementService {
	return &advertisementService{
		dao: dao.NewAdvertisementDao(dataSource.InstanceMaster()),
		}
}

func (s *advertisementService) Add(ad *models.Advertisement) error  {
	return s.dao.Add(ad)
}

func (s *advertisementService) GetAdvertisementList() *[]models.Advertisement  {
	return s.dao.GetAdvertisementList()
}

func (s *advertisementService) Update(newAd *models.Advertisement) error  {
	return s.dao.Update(newAd)
}

func (s *advertisementService) Delete(id int) error  {
	return s.dao.Delete(id)
}






