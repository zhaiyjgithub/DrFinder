package controllers

import (
	"DrFinder/src/models"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"DrFinder/src/utils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type UserTrackController struct {
	Ctx iris.Context
	UserTrackService service.UserTrackService
}

func (c *UserTrackController) BeforeActivation(b mvc.BeforeActivation)   {
	b.Handle(iris.MethodPost, utils.APIAddEvent, "AddTrackEvent")
}

func (c *UserTrackController) AddTrackEvent()  {
	type Param struct {
		Actions []models.UserAction
		Views []models.UserView
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	if len(param.Actions) != 0 {
		err  = c.UserTrackService.AddManyActionEvent(param.Actions)
		if err != nil  {
			response.Fail(c.Ctx, response.Err, err.Error(), nil)
			return
		}
	}

	if len(param.Views) != 0 {
		err = c.UserTrackService.AddManyViewTimeEvent(param.Views)
		if err != nil {
			response.Fail(c.Ctx, response.Err, err.Error(), nil)
			return
		}
	}

	response.Success(c.Ctx, response.Successful, nil)
}
