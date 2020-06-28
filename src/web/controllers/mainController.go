package controllers

import (
	"DrFinder/src/conf"
	"DrFinder/src/response"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"gopkg.in/go-playground/validator.v9"
)

var validate = validator.New()
var  j = jwt.New(jwt.Config{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
	return conf.JWRTSecret, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
	ErrorHandler: func(ctx iris.Context, e error) {
		response.Fail(ctx, response.Expire, e.Error(), nil)
	},
})