package controllers

import (
	"DrFinder/src/conf"
	"DrFinder/src/utils"
	"DrFinder/src/models"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"time"
)

type AnswerController struct {
	Ctx iris.Context
	Service service.AnswerService
	UserService service.UserService
}

func (c *AnswerController) BeforeActivation(b mvc.BeforeActivation)  {
	j := jwt.New(jwt.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return conf.JWRTSecret, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler: func(ctx iris.Context, e error) {
			response.Fail(ctx, response.Expire, e.Error(), nil)
		},
	})

	b.Handle(iris.MethodPost, utils.AddAnswer, "AddAnswer", j.Serve)
	b.Handle(iris.MethodPost, utils.DeleteDoctorById, "DeleteById", j.Serve)
	b.Handle(iris.MethodPost, utils.AddAnswerLikes, "AddAnswerLikes", j.Serve)
	b.Handle(iris.MethodPost, utils.GetAnswerListByPage, "GetAnswerListByPage", j.Serve)
}

func (c *AnswerController) AddAnswer()  {
	type Param struct {
		UserID      int  `validate:"gt=0"`
		PostID      int  `validate:"gt=0"`
		Description string `validate:"gt=0"`
	}

	var param Param

	err := utils.ValidateParam(c.Ctx, validate, &param)

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

	err := utils.ValidateParam(c.Ctx, validate, &param)

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

func (c *AnswerController) AddAnswerLikes()  {
	type Param struct {
		AnswerID int `validate:"gt=0"`
		UserID int `validate:"gt=0"`
	}

	var param Param

	err := utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	err = c.Service.AddLikes(param.AnswerID)

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

	err := utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	type AnswerInfo struct {
		ID int
		UserIcon string
		UserID int
		UserName string
		AnswerDate time.Time
		Description string
		PostID int
		Likes int
	}

	answerInfos := make([]AnswerInfo, 0, 0)
	answers := c.Service.GetAnswerListByPage(param.PostID, param.Page, param.PageSize)

	for _, answer := range answers {
		var answerInfo AnswerInfo
		answerInfo.ID = answer.ID
		answerInfo.UserID = answer.UserID
		answerInfo.AnswerDate = answer.CreatedAt
		answerInfo.Description = answer.Description
		answerInfo.PostID = answer.PostID
		answerInfo.Likes = answer.Likes

		userIno := c.UserService.GetUserById(answer.UserID)

		if userIno != nil {
			answerInfo.UserName = userIno.FirstName
			answerInfo.UserIcon = userIno.HeaderIcon
		}

		answerInfos = append(answerInfos, answerInfo)
	}

	response.Success(c.Ctx, response.Successful, answerInfos)
}
