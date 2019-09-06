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

type PostController struct {
	Ctx iris.Context
	Service service.PostService
}

func (c *PostController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle(iris.MethodPost, Utils.CreatePost, "CreatePost")
	b.Handle(iris.MethodPost, Utils.UpdatePost, "UpdatePost")
	b.Handle(iris.MethodPost,Utils.AddLikes, "AddLikes")
	b.Handle(iris.MethodPost, Utils.AddFavor, "AddFavor")
	b.Handle(iris.MethodPost, Utils.DeletePost, "DeletePost")
}

func (c *PostController) CreatePost() {
	type Param struct {
		UserID int `validate:"gt=0"`
		Type  int
		Title string `validate:"gt=0"`
		Description string `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	var post models.Post
	post.Type = param.Type
	post.UserID = param.UserID
	post.Title = param.Title
	post.Description = param.Description
	post.Favorites = 0
	post.Likes = 0
	post.Priority = 0
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	err = c.Service.Add(&post)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "create post fail", nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *PostController) UpdatePost()  {
	type Param struct {
		ID int `validate:"gt=0"`
		Description string `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	var post models.Post
	post.ID = param.ID
	post.Description = param.Description

	err = c.Service.Update(&post)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "update post faiel", nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *PostController) AddLikes()  {
	type Param struct {
		ID int `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	err = c.Service.AddLike(param.ID)

	if err != nil {
		response.Fail(c.Ctx, response.Err, response.NotFound, nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *PostController) AddFavor()  {
	type Param struct {
		ID int `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	err = c.Service.AddFavor(param.ID)

	if err != nil {
		response.Fail(c.Ctx, response.Err, response.NotFound, nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *PostController) DeletePost()  {
	type Param struct{
		ID int `validate:"gt=0"`
		UserID int `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	err = c.Service.DeleteByUser(param.ID, param.UserID)

	if err != nil {
		response.Fail(c.Ctx, response.Err, response.NotFound, nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}