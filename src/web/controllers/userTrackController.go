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
		Action *models.UserAction
		View *models.UserView
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, userValidate, &param)
	if err != nil {
		return
	}

	if param.Action == nil {
		response.Fail(c.Ctx, response.Err, err.Error(), nil)
	}else {
		err = c.UserTrackService.AddActionEvent(param.Action)
		if err != nil {
			response.Fail(c.Ctx, response.Err, err.Error(), nil)
		}else {
			response.Success(c.Ctx, response.Successful, nil)
		}
	}
}
