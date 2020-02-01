package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"gopkg.in/go-playground/validator.v9"
)

type AdminController struct {
	Ctx iris.Context
}

var adminValidate *validator.Validate

func (c *AdminController) BeforeActivation(b mvc.BeforeActivation)  {
	//
}