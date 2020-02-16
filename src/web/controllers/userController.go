package controllers

import (
	"DrFinder/src/utils"
	"DrFinder/src/models"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"gopkg.in/go-playground/validator.v9"
)

type UserController struct {
	Ctx iris.Context
	UserService service.UserService
}

var userValidate *validator.Validate

func (c *UserController) BeforeActivation(b mvc.BeforeActivation)  {
	userValidate = validator.New()

	b.Handle(iris.MethodPost, utils.CreateUser, "CreateUser")
	b.Handle(iris.MethodPost, utils.UpdatePassword, "UpdatePassword")
	b.Handle(iris.MethodPost, utils.UpdateUserInfo, "UpdateUserInfo")
	b.Handle(iris.MethodPost, utils.GetUserInfo, "GetUserInfo")
}

func (c *UserController) CreateUser() {
	type Param struct {
		Email      string `validate:"email"`
		Password   string `validate:"min=8,max=20"`
	}

	var param Param

	err := utils.ValidateParam(c.Ctx, userValidate, &param)

	if err != nil {
		return
	}

	var user models.User
	user.Email = param.Email
	user.Password = param.Password

	err = c.UserService.CreateUser(&user)

	if err != nil {
		response.Fail(c.Ctx, response.Err, err.Error(), nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *UserController) UpdatePassword() {
	type Param struct {
		Email string `validate:"email"`
		OldPwd string `validate:"min=6,max=20"`
		NewPwd string `validate:"min=6,max=20"`
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, userValidate, &param)

	if err != nil {
		return
	}

	err = c.UserService.UpdatePassword(param.Email, param.OldPwd, param.NewPwd)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "email or old password is wrong", nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *UserController) UpdateUserInfo()  {
	type Param struct {
		UserID int `validate:"gt=0"`
		Name string `validate:"gt=5"` //name length >= 6
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, userValidate, &param)
	if err != nil {
		return
	}

	var user models.User
	user.ID = param.UserID
	user.Name = param.Name

	err = c.UserService.UpdateUser(&user)
	if err != nil {
		response.Fail(c.Ctx, response.Err, "", nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *UserController) GetUserInfo() {
	type Param struct {
		UserID int `validate:"gt=0"`
	}
	
	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}
	
	user := c.UserService.GetUserById(param.UserID)
	response.Success(c.Ctx, response.Successful, *user)
}

