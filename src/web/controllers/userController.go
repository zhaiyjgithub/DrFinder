package controllers

import (
	"DrFinder/src/Utils"
	"DrFinder/src/models"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"gopkg.in/go-playground/validator.v9"
)

type UserController struct {
	Ctx iris.Context
	service service.UserService
}

var userValidate *validator.Validate

func (c *UserController) BeforeActivation(b mvc.BeforeActivation)  {
	userValidate = validator.New()

	b.Handle(iris.MethodPost, Utils.CreateUser, "CreateUser")
}

func (c *UserController) CreateUser() {
	type Param struct {
		Email      string `validate:"email"`
		Password   string `validate:"min=6,max=20"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, userValidate, &param)

	if err != nil {
		return
	}

	var user models.User
	user.Email = param.Email
	user.Password = param.Password

	err = c.service.CreateUser(&user)

	if err != nil {
		response.Fail(c.Ctx, response.Err, err.Error(), nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}