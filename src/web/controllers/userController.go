package controllers

import (
	"DrFinder/src/Utils"
	"DrFinder/src/conf"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

type UserController struct {
	Ctx iris.Context
	Service service.UserService
}

var userValidate *validator.Validate

func (c *UserController) BeforeActivation(b mvc.BeforeActivation)  {
	userValidate = validator.New()

	b.Handle(iris.MethodPost, Utils.CreateUser, "CreateUser")
	b.Handle(iris.MethodPost, Utils.UpdatePassword, "UpdatePassword")
	b.Handle(iris.MethodPost,Utils.ForgetPassword, "ForgetPassword")
	b.Handle(iris.MethodPost, Utils.UpdateUserInfo, "UpdateUserInfo")
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

	err = c.Service.CreateUser(&user)

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
	err := Utils.ValidateParam(c.Ctx, userValidate, &param)

	if err != nil {
		return
	}

	err = c.Service.UpdatePassword(param.Email, param.OldPwd, param.NewPwd)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "email or old password is wrong", nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *UserController) ForgetPassword()  {
	type Param struct {
		Email string `validate:"email"`
		Code string `validate:"len=6"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, userValidate, &param)
	if err != nil {
		return
	}

	_, err = c.Service.GetUserByEmail(param.Email)

	if err != nil {
		response.Fail(c.Ctx, response.NotExist, "email not is exist", nil)
		return
	}

	var code = VerificationCode{}

	cb := dataSource.Get(param.Email)

	if cb == nil {
		response.Fail(c.Ctx, response.Err, "verification code is invalidate", nil)
		return
	}

	err = json.Unmarshal(cb, &code)

	t, err := time.Parse(conf.TimeFormat, code.DateTime)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "parse time error", nil)
		return
	}

	tstamp := t.Unix()
	t2, err := time.Parse(conf.TimeFormat, time.Now().Format(conf.TimeFormat))
	nstamp := t2.Unix()

	if nstamp - tstamp > 60*10 {
		response.Fail(c.Ctx, response.Err, "verification code is invalidate", nil)
	}else if param.Code != code.Value {
		response.Fail(c.Ctx, response.Err, "verification code is invalidate", nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *UserController) UpdateUserInfo()  {
	type Param struct {
		ID int `validate:"gt=0"`
		LastName   string `validate:"gt=0"`
		FirstName  string `validate:"gt=0"`
		MiddleName string `validate:"gt=0"`
		Bio        string `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, userValidate, &param)

	if err != nil {
		return
	}

	var user models.User
	user.ID = param.ID
	user.LastName = param.LastName
	user.FirstName = param.FirstName
	user.MiddleName = param.MiddleName
	user.Bio = param.Bio

	err = c.Service.UpdateUser(&user)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "", nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}

}

