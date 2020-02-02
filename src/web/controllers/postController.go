package controllers

import (
	"DrFinder/src/Utils"
	"DrFinder/src/models"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"io/ioutil"
	"path/filepath"
	"time"
)

type PostController struct {
	Ctx iris.Context
	Service service.PostService
	AnswerService service.AnswerService
}

func (c *PostController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle(iris.MethodPost, Utils.CreatePost, "CreatePost")
	b.Handle(iris.MethodPost, Utils.UpdatePost, "UpdatePost")
	b.Handle(iris.MethodPost,Utils.AddLikes, "AddLikes")
	b.Handle(iris.MethodPost, Utils.AddFavor, "AddFavor")
	b.Handle(iris.MethodPost, Utils.DeletePost, "DeletePost")
	b.Handle(iris.MethodPost, Utils.GetPostByPage, "GetPostByPage")
	b.Handle(iris.MethodGet, Utils.ImgPost, "ImgPost")
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
		response.Fail(c.Ctx, response.Err, "update post fail", nil)
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

func (c *PostController) GetPostByPage()  {
	type Param struct {
		Type int
		Page int `validate:"gt=0"`
		PageSize int `validate:"gt=0"`
	}

	var param Param

	err := Utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	type PostInfo struct {
		post *models.Post
		answer *models.Answer
	}

	posts := c.Service.GetPostListByPage(param.Type, param.Page, param.PageSize)

	var postInfos []PostInfo
	for i := 0; i < len(posts); i ++  {
		post := &posts[i]
		answer := c.AnswerService.GetLastAnswer(post.ID)

		postInfos = append(postInfos, PostInfo{post:post, answer:answer})
	}

	fmt.Println(postInfos)
	response.Success(c.Ctx, response.Successful, posts)
}

func (c *PostController) ImgPost()  {
	fileName := c.Ctx.URLParam("name")

	filePath := fmt.Sprintf("./src/upload/post" + fileName)
	_ = c.Ctx.ServeFile(filePath, true)
}

func (c *PostController) CreatePost()  {
	type Param struct { //`validate:"gt=0"`
		UserID int `validate:"gt=0"`
		Type int `validate:"numeric"`
		Title string `validate:"gt=0"`
		Description string `validate:"gt=0"`
		Files []struct{
			Ext string `validate:"gt=0"`
			Base64Data string `validate:"base64"`
		}
	}

	var param Param
	err := Utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	files := param.Files
	failure := 0
	for _, file := range files {
		fileName := fmt.Sprintf("%s.%s", generateFileName(param.UserID), file.Ext)
		_, err = saveFile(file.Base64Data, "./src/upload/post", fileName)
		if err != nil {
			failure ++
		}
	}

	if failure == 0 {
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
			response.Fail(c.Ctx, response.Err, err.Error(), nil)
			return
		}

		response.Success(c.Ctx, response.Successful, nil)
	}else {
		response.Fail(c.Ctx, response.Err, "parse file failed", nil)
	}
}

func saveFile(imgBase64 string, destDir string, fileName string) (int64, error)  {
	imgBuffer, err  := base64.StdEncoding.DecodeString(imgBase64)

	if err != nil {
		fmt.Println("decode img base64 error")

		return 1, err
	}

	filePath := filepath.Join(destDir, fileName)
	err = ioutil.WriteFile(filePath, imgBuffer, 0666)

	if err != nil {
		return 1, err
	}

	return 0, nil
}

func generateFileName(userId int) string {
	data := []byte(fmt.Sprintf("%d-%d", userId, time.Now().Unix()))
	md5er := md5.New()
	md5er.Write(data)

	return hex.EncodeToString(md5er.Sum(nil))
}