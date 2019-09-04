package dao

import (
	"DrFinder/src/models"
	"errors"
	"github.com/jinzhu/gorm"
)

type UserDao struct {
	engine *gorm.DB
}

func NewUserDao(engine *gorm.DB) *UserDao {
	return &UserDao{engine: engine}
}

func (d *UserDao) CreateUser(user *models.User) error  {
	db := d.engine.Create(user)
	return db.Error
}

func (d *UserDao) GetUserById(id int) *models.User  {
	var user models.User

	d.engine.Where("id = ?", id).First(&user)
	return &user
}

func (d *UserDao) GetUserByEmail(email string) *models.User {
	var user models.User
	db := d.engine.Where("email = ?", email).First(&user)
	if db.Error != nil {
		return nil
	}

	return &user
}

func (d *UserDao) UpdateUser(newUser *models.User) error {
	var user models.User

	db := d.engine.Where("id = ?", newUser.ID).First(&user)

	if db.Error != nil {
		return db.Error
	}

	db = d.engine.Model(&user).Update(newUser)

	return db.Error
}

func (d *UserDao) DeleteUserById(id int) error {
	if id == 0 {
		return  errors.New("id must be > 0")
	}

	var user models.User
	user.ID = id

	db := d.engine.Delete(&user)

	return db.Error
}

func (d *UserDao) UpdatePassword(email string, oldPwd string, newPwd string) error {
	var user models.User

	db := d.engine.Where("email = ? AND password = ?", email, oldPwd).First(&user)

	if db.Error != nil {
		return db.Error
	}

	user.Password = newPwd
	db = d.engine.Model(&user).Update("password", newPwd)

	return db.Error
}