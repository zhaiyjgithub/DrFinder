package controllers

import (
	"DrFinder/src/Utils"
	"DrFinder/src/models"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type AnswerController struct {
	Ctx iris.Context
	Service service.AnswerService
}

func (c *AnswerController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle(iris.MethodPost, Utils.AddAnswer, "AddAnswer")
	b.Handle(iris.MethodPost, Utils.DeleteDoctorById, "DeleteById")
	b.Handle(iris.MethodPost, Utils.AddAnswerLikes, "AddLikes")
	b.Handle(iris.MethodPost, Utils.GetAnswerListByPage, "GetAnswerListByPage")
}

func (c *AnswerController) AddAnswer()  {
	type Param struct {
		UserID      int  `validate:"gt=0"`
		PostID      int  `validate:"gt=0"`
		Description string `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	var answer models.Answer
	answer.UserID = param.UserID
	answer.PostID = param.PostID
	answer.Description = param.Description
	answer.Likes = 0

	err = c.Service.AddAnswer(&answer)

	if err != nil {
		response.Fail(c.Ctx, response.Err, err.Error(), nil)
	}else {
		response.Success(c.Ctx, response.Successful,  nil)
	}
}

func (c *AnswerController) DeleteById()  {
	type Param struct {
		ID int `validate:"gt=0"`
		UserID int `validate:"get=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	err = c.Service.DeleteByUser(param.ID, param.UserID)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "", nil)
	}else {
		response.Success(c.Ctx, response.Successful,  nil)
	}
}

func (c *AnswerController) AddLikes()  {
	type Param struct {
		ID int `validate:"gt=0"`
		UserID int `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	err = c.Service.AddLikes(param.ID)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "", nil)
	}else {
		response.Success(c.Ctx, response.Successful,  nil)
	}
}

func (c *AnswerController) GetAnswerListByPage()  {
	type Param struct{
		UserID int `validate:"gt=0"`
		PostID int `validate:"gt=0"`
		Page int `validate:"gt=0"`
		PageSize int `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	answers := c.Service.GetAnswerListByPage(param.PostID, param.Page, param.PageSize)

	response.Success(c.Ctx, response.Successful, answers)
}
