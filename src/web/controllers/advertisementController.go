package controllers

import (
	"DrFinder/src/service"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type AdvertisementController struct {
	Ctx iris.Context
	Service service.AdvertisementService
}

func (c *AdvertisementController) BeforeActivation(b mvc.BeforeActivation)  {

}

