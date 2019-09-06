package controllers

import (
	"DrFinder/src/Utils"
	"DrFinder/src/models"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"time"
)

type AnswerController struct {
	Ctx iris.Context
	Service service.AnswerService
}

func (c *AnswerController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle(iris.MethodPost, Utils.AddAnswer, "AddAnswer")
	b.Handle(iris.MethodPost, Utils.DeleteDoctorById, "DeleteById")
	b.Handle(iris.MethodPost, Utils.AddAnswerLikes, "AddLikes")
}

func (c *AnswerController) AddAnswer()  {
	type Param struct {
		UserID      int  `validate:"gt=0"`
		Description string `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	var answer models.Answer
	answer.UserID = param.UserID
	answer.Description = param.Description
	answer.Likes = 0
	answer.CreatedAt = time.Now()
	answer.UpdatedAt = time.Now()

	err = c.Service.AddAnswer(&answer)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "", nil)
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
