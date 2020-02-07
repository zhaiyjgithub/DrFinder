package service

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	UpdatePassword(email string, oldPwd string, newPwd string) error
	UpdateUser(newUser *models.User) error
	GetUserById(id int) *models.User
	ResetPassword(email string, newPwd string) error
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

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.dao.GetUserByEmail(email)
}

func (s *userService) UpdatePassword(email string, oldPwd string, newPwd string) error  {
	return s.dao.UpdatePassword(email, oldPwd, newPwd)
}

func (s *userService) UpdateUser(newUser *models.User) error  {
	return s.dao.UpdateUser(newUser)
}

func (s *userService) GetUserById(id int) *models.User  {
	return s.dao.GetUserById(id)
}

func (s *userService) ResetPassword(email string, newPwd string) error  {
	return s.dao.ResetPassword(email, newPwd)
}