package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type UserService interface {
	CreateUser(user *models.User) error
}

type userService struct {
	dao *dao.UserDao
}

func NewUserService() UserService {
	return &userService{dao: dao.NewUserDao(dataSource.InstanceMaster())}
}

func (s *userService) CreateUser(user *models.User) error {
	return s.dao.CreateUser(user)
}