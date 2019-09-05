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